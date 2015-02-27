package records

import (
	"database/sql"
	//"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
	"strconv"
	"strings"
)

type sqliteRecordService struct {
	Conn *sql.DB
}

/*
This function takes a connection to a sqlite service, and returns a Record Service.  Should only be used with application setup code
*/
func NewSqliteService(conn *sql.DB) RecordService {
	return sqliteRecordService{Conn: conn}
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

		statement += "\n\t" + field.Name + " " + kind
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

func (service sqliteRecordService) Create(record interface{}) error {
	statement := service.CreateStatement(record)
	_, err := service.Conn.Exec(statement)
	return err
}

func (service sqliteRecordService) Insert(record interface{}) error {
	typ, val := typeAndVal(record)

	fields := goflect.GetInfo(record)
	statement := ""
	statement += "INSERT INTO " + typ.Name() + "("
	columns := make([]string, 0)
	for _, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		columns = append(columns, field.Name)
	}
	statement += strings.Join(columns, ", ")
	statement += " ) VALUES ("

	columns = make([]string, 0)
	for _, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		fieldVal := val.FieldByName(field.Name)
		columns = append(columns, wrap(fieldVal, field))
	}
	statement += strings.Join(columns, ", ")
	statement += " )"
	_, err := service.Conn.Exec(statement)
	return err
}

func (service sqliteRecordService) Update(record interface{}) error {
	typ, val := typeAndVal(record)

	fields := goflect.GetInfo(record)
	statement := "UPDATE " + typ.Name() + " SET "
	columns := make([]string, 0)
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(fieldVal.Interface()))
		}
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

func (service sqliteRecordService) Delete(record interface{}) error {
	typ, val := typeAndVal(record)

	fields := goflect.GetInfo(record)
	statement := "DELETE FROM " + typ.Name()
	match := matcher.NewStructMatcher()
	for _, field := range fields {
		fieldVal := val.FieldByName(field.Name)
		if field.IsPrimary {
			match.AddField(field.Name, matcher.Eq(fieldVal.Interface()))
		}
	}

	printer := matcher.NewSqlitePrinter()
	result, err := printer.Print(match)
	if err != nil {
		return err
	}
	statement += " WHERE " + result
	_, err = service.Conn.Exec(statement)
	return err
}

func (service sqliteRecordService) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service sqliteRecordService) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	typ, _ := typeAndVal(record)

	fields := goflect.GetInfo(record)
	statement := "SELECT "
	columns := make([]string, 0)
	for _, field := range fields {
		columns = append(columns, field.Name)
	}
	statement += strings.Join(columns, " , ")
	statement += " FROM " + typ.Name()

	printer := matcher.NewSqlitePrinter()
	result, err := printer.Print(match)
	if err != nil {
		return nil, err
	}
	statement += " WHERE " + result

	rows, err := service.Conn.Query(statement)
	if err != nil {
		return nil, err
	}

	output := func(r interface{}) bool {
		return nextRow(rows, r)
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

func nextRow(rows *sql.Rows, record interface{}) bool {
	next := rows.Next()
	if next {
		fields := goflect.GetInfo(record)
		val := reflect.ValueOf(record)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		vals := make([]interface{}, len(fields))
		addrs := make([]interface{}, len(fields))
		for i, _ := range vals {
			addrs[i] = &vals[i]
		}
		rows.Scan(addrs...)
		for i, field := range fields {
			fieldVal := val.FieldByName(field.Name)
			switch field.Kind {
			case reflect.Bool:
				fieldVal.Set(reflect.ValueOf(vals[i].(int64) != 0))
			case reflect.Float64:
				fieldVal.Set(reflect.ValueOf(float64(vals[i].(float64))))
			case reflect.Float32:
				fieldVal.Set(reflect.ValueOf(float32(vals[i].(float64))))
			case reflect.Int:
				fieldVal.Set(reflect.ValueOf(int(vals[i].(int64))))
			case reflect.Int64:
				fieldVal.Set(reflect.ValueOf(vals[i]))
			case reflect.Int32:
				fieldVal.Set(reflect.ValueOf(int32(vals[i].(int64))))
			case reflect.Int16:
				fieldVal.Set(reflect.ValueOf(int16(vals[i].(int64))))
			case reflect.Int8:
				fieldVal.Set(reflect.ValueOf(int8(vals[i].(int64))))
			case reflect.Uint:
				fieldVal.Set(reflect.ValueOf(uint(vals[i].(int64))))
			case reflect.Uint64:
				fieldVal.Set(reflect.ValueOf(uint64(vals[i].(int64))))
			case reflect.Uint32:
				fieldVal.Set(reflect.ValueOf(uint32(vals[i].(int64))))
			case reflect.Uint16:
				fieldVal.Set(reflect.ValueOf(uint16(vals[i].(int64))))
			case reflect.Uint8:
				fieldVal.Set(reflect.ValueOf(uint8(vals[i].(int64))))
			default:
				fieldVal.Set(reflect.ValueOf(string(vals[i].([]uint8))))
			}
		}

	}
	return next
}

func (service sqliteRecordService) Restrict(match matcher.Matcher) (RecordService, error) {
	return view{match: match, delegate: service}, nil
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
