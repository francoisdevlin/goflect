package goflect

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ReflectValue reflect.StructField

/*
This is used to determine the Field Info using reflection on a structure.  It will use the field's name as the name, and the field's type's Kind as the Kind
*/
func (field ReflectValue) GetFieldInfo() (output FieldInfo) {
	output.Name = field.Name
	output.Kind = field.Type.Kind()
	return output
}

/*
This is used to generate a SqlInfo field using reflection.  There is a struct tag, sql, that stores interesting information about the field.  The following are valid entries for the tag

    primary - denotes a primary key
    autoincrement - denotes that the field will be autoincremented by the db
    immutable - denotes that the field will be immutable, so it can't be updated
    unique - denotes that a fields value will be unique
    not-null - denotes that a field cannot be null
    index - denotes that the field is indexed for performance
    nominal - denotes that the field is a name alias for a record

*/
func (field ReflectValue) GetFieldSqlInfo() (output SqlInfo) {
	tags := field.Tag.Get("sql")

	output.IsPrimary = strings.Contains(tags, "primary")
	output.IsAutoincrement = strings.Contains(tags, "autoincrement")
	output.IsImmutable = strings.Contains(tags, "immutable") || output.IsAutoincrement
	output.IsUnique = strings.Contains(tags, "unique") || strings.Contains(tags, "primary")
	output.IsNullable = !(strings.Contains(tags, "not-null") || output.IsUnique)
	output.IsIndexed = strings.Contains(tags, "index") || output.IsUnique
	output.IsNominal = strings.Contains(tags, "nominal")

	return output
}

func (field ReflectValue) GetFieldUiInfo() (output UiInfo) {
	output.Description = field.Tag.Get("desc")
	output.Default = field.Tag.Get("default")
	output.FieldOrder, _ = strconv.ParseInt(field.Tag.Get("order"), 0, 64)

	tags := field.Tag.Get("ui")
	output.Hidden = strings.Contains(tags, "hidden")

	return output
}

func (field ReflectValue) GetFieldValidatorInfo() (output ValidatorInfo) {
	output.IsRequired = strings.Contains(field.Tag.Get("valid"), "required")
	output.MaxValue = field.Tag.Get("valid-max")
	output.MinValue = field.Tag.Get("valid-min")
	output.MatchRegex = field.Tag.Get("valid-regex")
	r, _ := regexp.Compile(",")
	output.InEnum = r.Split(field.Tag.Get("valid-enum"), -1)
	output.IsRequired = strings.Contains(field.Tag.Get("valid"), "required")
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
		fieldInfo := ReflectValue(field)

		knownFields[field.Name] = hydrateField(i, fieldInfo)
	}

	for _, field := range fieldNames {
		output = append(output, knownFields[field])
	}

	return output
}
