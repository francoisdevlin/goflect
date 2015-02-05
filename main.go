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
	IsNominal       bool
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

type Nominal struct {
	Id   int64
	Name string
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
	output.IsNominal = strings.Contains(tags, "nominal")

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

/**
This is yet another take on the active record pattern.
**/
type RecordService interface {
	Insert(record interface{})
	ReadAll(record interface{}) func(record interface{}) bool
	ReadAllWhere(record interface{}, conditions map[string]interface{}) func(record interface{}) bool
	ReadAllNominal(record interface{}) func(record *Nominal) bool
	ReadAllNominalWhere(record interface{}, conditions map[string]interface{}) func(record *Nominal) bool
	Get(id int64, record interface{})
	GetNominal(id int64) (output Nominal)
	GetByNominal(name string, record interface{})
	GetNominalByNominal(name string) (output Nominal)
	Update(record interface{})
	Delete(record interface{})
	DeleteById(Id int64)
}

type SqliteRecordService struct {
	Conn *sql.DB
}

func (service SqliteRecordService) Insert(record interface{}) {
	message := InsertSQLiteRecord(record)
	_, _ = service.Conn.Exec(message)
}

func (service SqliteRecordService) ReadAll(record interface{}) func(record interface{}) bool {
	conditions := make(map[string]interface{})
	return service.ReadAllWhere(record, conditions)
}

func (service SqliteRecordService) ReadAllWhere(record interface{}, conditions map[string]interface{}) func(record interface{}) bool {
	message := ListSQLiteRecordWhere(record, conditions)
	rows, _ := service.Conn.Query(message)

	output := func(r interface{}) bool {
		return NextRow(rows, r)
	}

	return output
}

func (service SqliteRecordService) ReadAllNominal(record interface{}) func(record *Nominal) bool {
	message := ListSQLiteNominal(record)
	rows, _ := service.Conn.Query(message)

	output := func(r *Nominal) bool {
		return NextRow(rows, r)
	}

	return output
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

func wrap(fieldVal reflect.Value, field Info) string {
	output := ""
	switch {
	case field.Kind == reflect.Int:
	case field.Kind == reflect.Int8:
	case field.Kind == reflect.Int16:
	case field.Kind == reflect.Int32:
	case field.Kind == reflect.Int64:
		output = "" + strconv.FormatInt(fieldVal.Int(), 10) + ""
	case true:
		output = "\"" + fieldVal.String() + "\""
	}
	return output
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
		statement += " " + wrap(fieldVal, field)
		if i != len(fields)-1 {
			statement += ","
		}
	}
	statement += " )"
	return statement
}

func ListSQLiteRecord(record interface{}) (statement string) {
	conditions := make(map[string]interface{})
	return ListSQLiteRecordWhere(record, conditions)
}
func ListSQLiteRecordWhere(record interface{}, conditions map[string]interface{}) (statement string) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	segments := make([]string, 0)

	fields := GetInfo(record)
	statement = ""
	statement += "SELECT "
	for i, field := range fields {
		statement += " " + field.Name
		if i != len(fields)-1 {
			statement += ","
		}
		if conditional, present := conditions[field.Name]; present {
			condVal := reflect.ValueOf(conditional)
			if condVal.Kind() == reflect.Ptr {
				condVal = condVal.Elem()
			}
			segments = append(segments, fmt.Sprintf("%v = %v", field.Name, wrap(condVal, field)))
		}
	}
	statement += " FROM " + typ.Name()

	if len(segments) > 0 {
		statement += " WHERE " + strings.Join(segments, " AND ")
	}

	return statement
}

func ListSQLiteNominal(record interface{}) (statement string) {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := GetInfo(record)
	id, nominal := "", ""
	for _, field := range fields {
		if field.IsPrimary {
			id = field.Name
		}
		if field.IsNominal {
			nominal = field.Name
		}
	}
	statement = fmt.Sprintf("SELECT %v, %v FROM %v", id, nominal, typ.Name())

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
