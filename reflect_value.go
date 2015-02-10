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

//These don't really go here...
func hydrateField(i int, field FieldDescription) Info {
	output := Info{
		FieldInfo:     field.GetFieldInfo(),
		SqlInfo:       field.GetFieldSqlInfo(),
		UiInfo:        field.GetFieldUiInfo(),
		ValidatorInfo: field.GetFieldValidatorInfo(),
	}
	if output.FieldOrder == 0 {
		output.FieldOrder = int64(i)
	}
	return output
}

func GetInfo(record interface{}) (output []Info) {
	typ := reflect.TypeOf(record)
	if res, err := record.(reflect.StructField); err {
		typ = res.Type
	}
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	knownFields := make(map[string]Info)
	fieldNames := make([]string, 0)
	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			anonFields := GetInfo(field)
			for _, anonField := range anonFields {
				//Add the new field if we don't know about it
				if _, present := knownFields[anonField.Name]; !present {
					fieldNames = append(fieldNames, anonField.Name)
					knownFields[anonField.Name] = anonField
				}
			}
			continue
		}
		if _, present := knownFields[field.Name]; !present {
			fieldNames = append(fieldNames, field.Name)
		}
		fieldInfo := reflectValue(field)

		knownFields[field.Name] = hydrateField(i, fieldInfo)
	}

	for _, field := range fieldNames {
		output = append(output, knownFields[field])
	}

	return output
}
