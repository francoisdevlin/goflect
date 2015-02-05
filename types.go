package main

import (
	//"database/sql"
	//"fmt"
	"reflect"
	//"strconv"
	//"strings"
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
