package goflect

import (
	"reflect"
	"strconv"
	"strings"
)

type reflectValue reflect.StructField

/*
This is used to determine the Field Info using reflection on a structure.  It will use the field's name as the name, and the field's type's Kind as the Kind
*/
func (field reflectValue) GetFieldInfo() (output FieldInfo) {
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
func (field reflectValue) GetFieldSqlInfo() (output SqlInfo) {
	tags := field.Tag.Get(TAG_SQL)

	output.IsPrimary = strings.Contains(tags, SQL_PRIMARY)
	output.IsAutoincrement = strings.Contains(tags, SQL_AUTOINC)
	output.IsImmutable = strings.Contains(tags, SQL_IMMUTABLE) || output.IsAutoincrement
	output.IsUnique = strings.Contains(tags, SQL_UNIQUE) || output.IsPrimary
	output.IsNullable = !(strings.Contains(tags, SQL_NULLABLE) || output.IsUnique)
	output.IsIndexed = strings.Contains(tags, SQL_INDEX) || output.IsUnique
	output.IsNominal = strings.Contains(tags, SQL_NOMINAL)
	output.IsSqlIgnored = strings.Contains(tags, SQL_IGNORE)

	output.SqlColumn = field.Tag.Get(TAG_SQL_COLUMN)
	output.ChildOf = field.Tag.Get(TAG_SQL_CHILD)
	output.Extends = field.Tag.Get(TAG_SQL_EXTEND)

	return output
}

/*
This is used to generate a UiInfo field using reflection.  There are several struct tags that are used to store interesting information about a field

    desc - This stores a human readable description for a tooltip
    default - This stores the default value for the field.  Must be compatible with the type
    order - This controls the order for the field to appear in web forms

There is also a "flag tag", "ui", with the following entries possible

    hidden - This controls if the user can see the field
*/
func (field reflectValue) GetFieldUiInfo() (output UiInfo) {
	output.Description = field.Tag.Get(TAG_DESC)
	output.Default = field.Tag.Get(TAG_DEFAULT)
	output.FieldOrder, _ = strconv.ParseInt(field.Tag.Get(TAG_ORDER), 0, 64)

	tags := field.Tag.Get(TAG_UI)
	output.IsHidden = strings.Contains(tags, UI_HIDDEN)
	output.IsRedacted = strings.Contains(tags, UI_REDACTED)

	return output
}

/*
This is used to get information our of the valid tag, and into the data structure.
*/
func (field reflectValue) GetFieldValidatorInfo() (output ValidatorInfo) {
	output.ValidExpr = field.Tag.Get(TAG_VALID)
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
