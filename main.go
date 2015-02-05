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
// * Create
// * READ 1
// * READ WHERE (Excel autofilter)
// * Update 1
// * Delete 1
//MySQL table creation
//MySQL record maniputalor
//Standard REST endpoints
//"Nominal Dropdowns" from SQL backend
package main

import (
	//"database/sql"
	"fmt"
	"reflect"
	//"strconv"
	//"strings"
)

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

func main() {
	//GetRecords([]int{1, 2, 3}, &[]int{})
	//CreateSQLiteTable(&Foo{})
	fmt.Println("Hello")
}
