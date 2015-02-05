package main

import (
	//"database/sql"
	//"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
