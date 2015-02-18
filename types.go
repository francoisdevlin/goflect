package goflect

import (
	"reflect"
)

/*
The FieldInfo struct is used to store two pieces of information about the field, its name and Kind.
*/
type FieldInfo struct {
	Name string       `desc:"This is the name of the field in the struct.  It is authoritative" sql:"primary"`
	Kind reflect.Kind `desc:"This is the golang kind, from the reflect pacakge.  It controls dispatch"`
}

type SqlInfo struct {
	IsPrimary       bool `desc:"This indicates if the field is the primary key.  It will imply uniqueness, and all that follows"`
	IsAutoincrement bool `desc:"This indicates if the field is an integet auto increment.  It implies immutability on the field"`
	IsUnique        bool `desc:"This indicates if the field must be unique.  It implies not not and and index"`
	IsNullable      bool `desc:"This controls if a field is nullable."`
	IsIndexed       bool `desc:"This controls if a field is indexed."`
	IsNominal       bool `desc:"This indicates the nominal field for a type, which generates dropdowns.  There may be only one per structure"`
	IsImmutable     bool `desc:"This controls if a field is immutable.  It will make the field write once."`
}

type ValidatorInfo struct {
	IsRequired bool `desc:"This determines if a field is required"`
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

/*
This is the structure that holds all of the field information.  It is what is intended to be consumed by the reflection engine
*/
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
	DeleteById(id int64, record interface{})
}
