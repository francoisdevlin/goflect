package records

import (
	"database/sql"
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
	"strconv"
	"strings"
)

type sqliteRecordService struct {
	Conn *sql.DB
}

func typeAndVal(record interface{}) (reflect.Type, reflect.Value) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return typ, val
}

func sqliteLookupMap() map[reflect.Kind]string {
	lookup := map[reflect.Kind]string{
		reflect.Bool:    "integer",
		reflect.Int:     "integer",
		reflect.Int8:    "integer",
		reflect.Int16:   "integer",
		reflect.Int32:   "integer",
		reflect.Int64:   "integer",
		reflect.Uint:    "integer",
		reflect.Uint8:   "integer",
		reflect.Uint16:  "integer",
		reflect.Uint32:  "integer",
		reflect.Uint64:  "integer",
		reflect.Float32: "real",
		reflect.Float64: "real",
	}
	return lookup
}

/*
This renders the sql statement to create a table, based on the provided record
*/
func (service sqliteRecordService) CreateStatement(record interface{}) string {
	typ, _ := typeAndVal(record)

	fields := goflect.GetInfo(record)
	lookup := sqliteLookupMap()
	statement := ""
	statement += "CREATE TABLE IF NOT EXISTS " + typ.Name() + "("
	for i, field := range fields {
		kind, present := lookup[field.Kind]
		if !present {
			kind = "string"
		}

		statement += "\n\t`" + field.Name + "` " + kind
		if field.IsPrimary {
			statement += " primary key"
		}
		if field.IsAutoincrement {
			statement += " autoincrement"
		}
		if !field.IsNullable {
			statement += " not null"
		}
		if i != len(fields)-1 {
			statement += ","
		}
	}
	statement += "\n)"
	return statement
}

func (service sqliteRecordService) Define(record interface{}) error {
	statement := service.CreateStatement(record)
	_, err := service.Conn.Exec(statement)
	return err
}

func (service sqliteRecordService) createAll(record interface{}) error {
	typ, val := typeAndVal(record)

	typ = reflect.TypeOf(record)
	for typ.Kind() == reflect.Array || typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	fields := goflect.GetInfo(val.Index(0).Interface())
	statement := ""
	statement += "INSERT INTO `" + typ.Name() + "`("
	columns := make([]string, 0)
	for _, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		columns = append(columns, "`"+field.Name+"`")
	}
	statement += strings.Join(columns, ", ")
	statement += " ) VALUES "

	columns = make([]string, 0)
	for i := 0; i < val.Len(); i++ {
		columns = append(columns, uglyGuy(fields, val.Index(i).Interface()))
	}
	statement += strings.Join(columns, ", ")
	_, err := service.Conn.Exec(statement)
	return err
}

func uglyGuy(fields []goflect.Info, record interface{}) string {
	_, val := typeAndVal(record)
	columns := make([]string, 0)
	for _, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		fieldVal := val.FieldByName(field.Name)
		columns = append(columns, wrap(fieldVal, field))
	}
	statement := strings.Join(columns, ", ")
	statement = "( " + statement + " )"
	return statement
}

func (service sqliteRecordService) updateAll(record interface{}, match matcher.Matcher) error {
	typ, val := typeAndVal(record)

	fields := goflect.GetInfo(record)
	statement := "UPDATE " + typ.Name() + " SET "
	columns := make([]string, 0)
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		temp := field.Name + "=" + wrap(fieldVal, field)
		columns = append(columns, temp)
	}
	statement += strings.Join(columns, ", ")

	printer := matcher.NewSqlitePrinter()
	result, err := printer.Print(match)
	if err != nil {
		return err
	}
	statement += " WHERE " + result

	_, err = service.Conn.Exec(statement)
	return err
}

func (service sqliteRecordService) deleteAll(record interface{}, match matcher.Matcher) error {
	typ, _ := typeAndVal(record)

	statement := "DELETE FROM " + typ.Name()

	printer := matcher.NewSqlitePrinter()
	result, err := printer.Print(match)
	if err != nil {
		return err
	}
	statement += " WHERE " + result
	_, err = service.Conn.Exec(statement)
	return err
}

type Edge struct {
	A, B, Table string
}

func determineEdges(records []interface{}) (output []Edge) {
	primaryKeys := map[string]string{}
	for _, record := range records {
		typ, _ := typeAndVal(record)
		name := typ.Name()
		fields := goflect.GetInfo(record)
		for _, field := range fields {
			if field.IsPrimary {
				primaryKeys[name] = field.Name
			}
		}
	}
	for _, record := range records {
		typ, _ := typeAndVal(record)
		name := typ.Name()
		parent := ""
		localField := ""
		fields := goflect.GetInfo(record)
		for _, field := range fields {
			if field.ChildOf != "" {
				parent = field.ChildOf
				localField = field.Name
			}
			if field.Extends != "" {
				parent = field.Extends
				localField = field.Name
			}
		}
		if sourceField, hit := primaryKeys[parent]; hit {
			output = append(output, Edge{A: parent + "." + sourceField, B: name + "." + localField, Table: name})
		}
	}
	return output
}

func (service sqliteRecordService) readAll(query matcher.Matcher, records ...interface{}) (func(record ...interface{}) bool, error) {
	statement := "SELECT "
	//path := make(edge
	types := make([]reflect.Type, 0, 0)
	columns := make([]string, 0)
	for _, record := range records {
		typ, _ := typeAndVal(record)

		fields := goflect.GetInfo(record)
		for _, field := range fields {
			columns = append(columns, typ.Name()+".`"+field.Name+"`")
		}
		types = append(types, typ)

	}
	statement += strings.Join(columns, " , ")
	if len(types) == 1 {
		statement += " FROM " + types[0].Name()
	} else {
		statement += " FROM " + types[0].Name()
		edges := determineEdges(records)
		for _, edge := range edges {
			statement += fmt.Sprintf(" INNER JOIN %v ON %v = %v", edge.Table, edge.A, edge.B)
		}
	}

	printer := matcher.NewSqlitePrinter()
	result, err := printer.Print(query)
	if err != nil {
		return nil, err
	}
	statement += " WHERE " + result

	rows, err := service.Conn.Query(statement)
	if err != nil {
		//fmt.Println(statement)
		return nil, err
	}

	output := func(r ...interface{}) bool {
		return nextRow(rows, r...)
	}

	return output, nil
}

func wrap(fieldVal reflect.Value, field goflect.Info) string {
	output := ""
	switch fieldVal.Kind() {
	case reflect.Bool:
		if fieldVal.Bool() {
			output = "1"
		} else {
			output = "0"
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		output = "" + strconv.FormatInt(int64(fieldVal.Uint()), 10) + ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		output = "" + strconv.FormatInt(fieldVal.Int(), 10) + ""
	case reflect.Float32:
		output = "" + strconv.FormatFloat(fieldVal.Float(), 'f', 10, 32) + ""
	case reflect.Float64:
		output = "" + strconv.FormatFloat(fieldVal.Float(), 'f', 10, 64) + ""
	default:
		output = "\"" + fieldVal.String() + "\""
	}
	return output
}

func nextRow(rows *sql.Rows, records ...interface{}) bool {
	next := rows.Next()
	if next {
		total := 0
		coerce := func(v interface{}, fVal reflect.Value) {
			localVal := reflect.ValueOf(v)
			if fVal.Type() != localVal.Type() {
				localVal = localVal.Convert(fVal.Type())
			}
			fVal.Set(localVal)
		}

		for _, record := range records {
			fields := goflect.GetInfo(record)
			total += len(fields)
		}

		vals := make([]interface{}, total)
		addrs := make([]interface{}, total)
		for i, _ := range vals {
			addrs[i] = &vals[i]
		}

		rows.Scan(addrs...)

		offset := 0
		for _, record := range records {
			fields := goflect.GetInfo(record)
			val := reflect.ValueOf(record)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			for i, field := range fields {
				fieldVal := val.FieldByName(field.Name)
				switch field.Kind {
				case reflect.Bool:
					coerce(vals[offset+i].(int64) != 0, fieldVal)
				default:
					coerce(vals[offset+i], fieldVal)
				}
			}
			offset += len(fields)
		}

	}
	return next
}

func RailsConvention(record interface{}) func(interface{}) (interface{}, error) {
	typ, val := typeAndVal(record)

	closedRecords := make(map[string]interface{})
	closedName := typ.Name()
	primaryName := ""
	var closedId interface{} = nil
	fields := goflect.GetInfo(record)
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		if field.IsPrimary {
			closedId = fieldVal.Interface()
			primaryName = field.Name
			//closedRecords[typ.Name()] = fieldVal.Interface()
		} else if strings.Contains(field.Name, "Id") {
			closedRecords[field.Name] = fieldVal.Interface()
		}
	}

	return func(other interface{}) (interface{}, error) {
		typ, val := typeAndVal(other)
		//This will find a parent object.  E.g., a device given an object
		parentColumn := typ.Name() + "Id"
		if value, hit := closedRecords[parentColumn]; hit {
			otherField := val.FieldByName(primaryName)
			otherField.Set(reflect.ValueOf(value))
			return other, nil
		}
		fields := goflect.GetInfo(other)
		for _, field := range fields {
			if field.Name == closedName+"Id" {
				fieldVal := val.FieldByName(field.Name)
				fieldVal.Set(reflect.ValueOf(closedId))
				return other, nil
			}
		}
		return other, nil
	}

}
