package goflect

import (
	"reflect"
)

const (
	TAG_DESC       string = "desc"
	TAG_SQL               = "sql"
	TAG_SQL_COLUMN        = "sql-column"
	TAG_SQL_CHILD         = "sql-child"
	TAG_SQL_EXTEND        = "sql-extend"
	TAG_VALID             = "valid"
	TAG_DEFAULT           = "default"
	TAG_ORDER             = "order"
	TAG_UI                = "ui"
	TAG_UI_NAME           = "ui-name"
)

const (
	SQL_PRIMARY   string = "primary"
	SQL_AUTOINC          = "autoincrement"
	SQL_IMMUTABLE        = "immutable"
	SQL_UNIQUE           = "unique"
	SQL_NULLABLE         = "not-null"
	SQL_INDEX            = "index"
	SQL_NOMINAL          = "nominal"
	SQL_IGNORE           = "ignored"
)

const (
	UI_HIDDEN   string = "hidden"
	UI_REDACTED        = "redacted"
)

var (
	TAGS = []string{
		TAG_DESC,
		TAG_DEFAULT,
		TAG_VALID,
		TAG_SQL,
		TAG_SQL_COLUMN,
		TAG_SQL_CHILD,
		TAG_SQL_EXTEND,
		TAG_UI,
		TAG_UI_NAME,
		TAG_ORDER,
	}
	SQL_FIELDS = []string{
		SQL_PRIMARY,
		SQL_AUTOINC,
		SQL_UNIQUE,
		SQL_IMMUTABLE,
		SQL_NOMINAL,
		SQL_NULLABLE,
		SQL_INDEX,
		SQL_IGNORE,
	}
	UI_FIELDS = []string{
		UI_HIDDEN,
		UI_REDACTED,
	}
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
	IsSqlIgnored    bool `desc:"This controls if the field is ifnored entirely by all sql operations.  It will not render down to sql queries.  However, it may still be used by in memory transforms"`
	Bacon           string
	SqlColumn       string `desc:"This is the actual sql column to use.  Leaving it blank to allow the engine to determine the value based on the Name property"`
	ChildOf         string `desc:"This describes the child relationship that a record has.  It points to a type"`
	Extends         string `desc:"This describes the extension relationship that a table has"`
}

type ValidatorInfo struct {
	ValidExpr string `desc:"This stores a validation expresion that drives the default struct matcher"`
}

type UiInfo struct {
	Description string `desc:"This is the user facing description of a field"`
	FieldOrder  int64  `desc:"This controls the field display order"`
	IsHidden    bool   `desc:"This controls if the user can see the field at all"`
	IsRedacted  bool   `desc:"This controls if a field is redacted.  It will show up as stars in user input"`
	Default     string `desc:"This is the default value of the field in the UI"`
	DisplayName string `desc:"This is the display name shown to the user.  Leving it blank will allow the engine to determine the value used"`
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

type FieldDescription interface {
	GetFieldInfo() FieldInfo
	GetFieldSqlInfo() SqlInfo
	GetFieldUiInfo() UiInfo
	GetFieldValidatorInfo() ValidatorInfo
}
