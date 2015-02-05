//Things to generate
//Go type parser
//Type JSON Marshall/UnMarshall (DONE!  Thanks Go!)
//CLI interface (Action Framework)
//Web Component
//Ext Model
//Ext Store
//Ext Table
//Ext Default Form
//JSON Marshall/UnMarshall (leverage existing, quick hack)
//Protobuff Marshall/UnMarshall
//Avro Marshall/UnMarshall
//EDN Marshall/UnMarshall
//map[string]interaface{} reader
//Curses Interface
//SQLite table creation
//SQLite record maniputalor
//MySQL table creation
//MySQL record maniputalor
//Standard REST endpoints
//"Nominal Dropdowns" from SQL backend
package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type FieldInfo struct {
	Name string
	Kind reflect.Kind
}

type SqlInfo struct {
	IsPrimary       bool
	IsAutoincrement bool
	IsUnique        bool
	IsNullable      bool
	IsIndexed       bool
}

type ValidatorInfo struct {
	IsRequired bool
	MaxValue   string
	MinValue   string
	MatchRegex string
	InEnum     []string
}

type UiInfo struct {
	Description string
	FieldOrder  int64
	Hidden      bool //This controls if the value can ever be interacted with
	Default     string
}

type Info struct {
	FieldInfo
	SqlInfo
	ValidatorInfo
	UiInfo
}

type FieldDescription interface {
	GetFieldInfo() FieldInfo
	GetFieldSqlInfo() SqlInfo
	GetFieldUiInfo() UiInfo
	GetFieldValidatorInfo() ValidatorInfo
}

type reflectValue reflect.StructField

func (field reflectValue) GetFieldInfo() (output FieldInfo) {
	output.Name = field.Name
	output.Kind = field.Type.Kind()
	return output
}

func (field reflectValue) GetFieldSqlInfo() (output SqlInfo) {
	tags := field.Tag.Get("sql")

	output.IsPrimary = strings.Contains(tags, "primary")
	output.IsAutoincrement = strings.Contains(tags, "autoincrement")
	output.IsUnique = strings.Contains(tags, "unique") || strings.Contains(tags, "primary")
	output.IsNullable = !strings.Contains(tags, "not-null") || output.IsUnique
	output.IsIndexed = strings.Contains(tags, "index") || output.IsUnique

	return output
}

func (field reflectValue) GetFieldUiInfo() (output UiInfo) {
	output.Description = field.Tag.Get("desc")
	output.Default = field.Tag.Get("default")
	output.FieldOrder, _ = strconv.ParseInt(field.Tag.Get("order"), 0, 64)

	tags := field.Tag.Get("ui")
	output.Hidden = strings.Contains(tags, "hidden")

	return output
}

func (field reflectValue) GetFieldValidatorInfo() (output ValidatorInfo) {
	return output
}

func getInfo(fields []FieldDescription) (output []Info) {
	for i, field := range fields {
		temp := Info{
			FieldInfo:     field.GetFieldInfo(),
			SqlInfo:       field.GetFieldSqlInfo(),
			UiInfo:        field.GetFieldUiInfo(),
			ValidatorInfo: field.GetFieldValidatorInfo(),
		}
		if temp.FieldOrder == 0 {
			temp.FieldOrder = int64(i)
		}
		output = append(output, temp)
	}

	return output
}

func GetInfo(record interface{}) (output []Info) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// loop through the struct's fields and set the map
	fields := make([]FieldDescription, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fields = append(fields, reflectValue(field))
	}

	output = getInfo(fields)

	return output
}

type SqlHydration interface {
	GetRecords(ids []int)
	GetRecord(ids int)
	UpdateRecords()
	InsertRecords()
	DeleteRecords()
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

func CreateSQLiteTable(record interface{}) (statement string) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := GetInfo(record)
	lookup := sqliteLookupMap()
	statement = ""
	statement += "CREATE TABLE " + typ.Name() + "("
	for i, field := range fields {
		kind, present := lookup[field.Kind]
		if !present {
			kind = "string"
		}

		statement += "\n" + field.Name + " " + kind
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

func InsertSQLiteRecord(record interface{}) (statement string) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := GetInfo(record)
	statement = ""
	statement += "INSERT INTO " + typ.Name() + "("
	for i, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		statement += " " + field.Name
		if i != len(fields)-1 {
			statement += ","
		}
	}
	statement += " ) VALUES ("

	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i, field := range fields {
		if field.IsAutoincrement {
			continue
		}
		fieldVal := val.FieldByName(field.Name)
		statement += " \"" + fieldVal.String() + "\""
		if i != len(fields)-1 {
			statement += ","
		}
	}
	statement += " )"
	//lookup := sqliteLookupMap()
	return statement
}

func ListSQLiteRecord(record interface{}) (statement string) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := GetInfo(record)
	statement = ""
	statement += "SELECT "
	for i, field := range fields {
		statement += " " + field.Name
		if i != len(fields)-1 {
			statement += ","
		}
	}
	statement += " FROM " + typ.Name()

	return statement
}

func NextRow(rows *sql.Rows, record interface{}) bool {
	next := rows.Next()
	if next {
		fields := GetInfo(record)
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
			if field.Kind == reflect.String {
				fieldVal.Set(reflect.ValueOf(string(vals[i].([]uint8))))
			} else {
				fieldVal.Set(reflect.ValueOf(vals[i]))
			}
		}

	}
	return next
}

func main() {
	//GetRecords([]int{1, 2, 3}, &[]int{})
	//CreateSQLiteTable(&Foo{})
	fmt.Println("Hello")
}
