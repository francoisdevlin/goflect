package goflect

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git"
	"go/token"
	"reflect"
	"regexp"
)

type ErrorCode int

const (
	NOMINAL_MISMATCH ErrorCode = iota
	TAG_ERROR
	NOMINAL_MISCOUNT
	PRIMARY_MISCOUNT
)

type ValidationError struct {
	Code    ErrorCode
	Message string
}

func (v ValidationError) Error() string {
	return v.Message
}

type Result struct {
	Error    ValidationError
	Position token.Position
}

func ValidateType(record interface{}, list StructList) []Result {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	//fieldChecks := make([]func(f goflect.Info) []error, 0)
	fieldChecks := []func(f goflect.Info) []error{
		Nominal,
	}

	//recordChecks := make([]func(f []goflect.Info) []error, 0)
	recordChecks := []func(f []goflect.Info) []error{
		NominalOnce,
		PrimaryOnce,
	}

	output := make([]Result, 0)
	fields := goflect.GetInfo(record)

	for _, recordCheck := range recordChecks {
		errors := recordCheck(fields)
		for _, err := range errors {
			cast, _ := err.(ValidationError)
			output = append(output, Result{Error: cast, Position: list.Position})
		}
	}
	for _, field := range fields {
		pos := list.Position
		if fieldStruct, present := list.Structs[field.Name]; present {
			pos = fieldStruct.Position
		}
		valType, _ := typ.FieldByName(field.Name)
		errors := StructTag(string(valType.Tag))
		for _, err := range errors {
			cast, _ := err.(ValidationError)
			output = append(output, Result{Error: cast, Position: pos})
		}
		for _, fieldCheck := range fieldChecks {
			errors := fieldCheck(field)
			for _, err := range errors {
				cast, _ := err.(ValidationError)
				output = append(output, Result{Error: cast, Position: pos})
			}
		}
	}
	return output
}

/*
This test ensures that there is only one field marked nominal for a given type
*/
func NominalOnce(fields []goflect.Info) []error {
	nominalFields := make([]string, 0)
	for _, field := range fields {
		if field.IsNominal {
			nominalFields = append(nominalFields, field.Name)
		}
	}
	errors := make([]error, 0)
	if len(nominalFields) > 1 {
		errors = append(errors, ValidationError{
			Code:    NOMINAL_MISCOUNT,
			Message: fmt.Sprintf("There can be only one nominal field, but the following are marked, %v", nominalFields),
		})
	}
	return errors
}

/*
This test ensures that there is only one field marked primary for a given type
*/
func PrimaryOnce(fields []goflect.Info) []error {
	primaryFields := make([]string, 0)
	for _, field := range fields {
		if field.IsPrimary {
			primaryFields = append(primaryFields, field.Name)
		}
	}
	errors := make([]error, 0)
	if len(primaryFields) > 1 {
		errors = append(errors, ValidationError{
			Code:    PRIMARY_MISCOUNT,
			Message: fmt.Sprintf("There can be only one primary field, but the following are marked, %v", primaryFields),
		})
	}
	return errors
}

/*
This checks the conditions around a nominal field.  The requirements are that the nominal field is type string, and that the nominal field is unique
*/
func Nominal(f goflect.Info) []error {
	errors := make([]error, 0)
	if f.IsNominal && f.Kind != reflect.String {
		errors = append(errors, ValidationError{
			Code:    NOMINAL_MISMATCH,
			Message: fmt.Sprintf("Field %v is marked nominal, but is kind %v", f.Name, f.Kind),
		})
	}
	if f.IsNominal && !f.IsUnique {
		errors = append(errors, ValidationError{
			Code:    NOMINAL_MISMATCH,
			Message: fmt.Sprintf("Field %v is marked nominal, but is not unique", f.Name),
		})
	}
	return errors
}

/*
This checks the syntax of the contents of the struct tag, so that the reflection engine works properly.  It will ensure that each tag is unique, the tag can be parsed, and that the tag is followed by a quoted string
*/
func StructTag(message string) []error {
	whitespace, _ := regexp.Compile("^[\\s,]+")
	symbol, _ := regexp.Compile("^[a-zA-Z_]\\w*")
	operators, _ := regexp.Compile("^:\"")
	quote, _ := regexp.Compile("^\"(?:\\\\?.)*?\"")
	errors := make([]error, 0)
	tokens := make([]string, 0)
	for len(message) > 0 {
		switch {
		case whitespace.MatchString(message):
			message = whitespace.ReplaceAllString(message, "")
		case operators.MatchString(message):
			message = operators.ReplaceAllString(message, "\"")
		case symbol.MatchString(message):
			tokens = append(tokens, symbol.FindString(message))
			message = symbol.ReplaceAllString(message, "")
		case quote.MatchString(message):
			tokens = append(tokens, quote.FindString(message))
			message = quote.ReplaceAllString(message, "")
		default:
			errors = append(errors, ValidationError{Code: TAG_ERROR, Message: message})
			message = ""
		}
	}

	if len(errors) > 0 {
		return errors
	}
	if len(tokens)%2 != 0 {
		errors = append(errors, ValidationError{Code: TAG_ERROR, Message: "There are the wrong number of tokens present"})
	}
	if len(errors) > 0 {
		return errors
	}

	iteration := 0
	tagKeys := map[string]string{}
	for iteration < len(tokens) {
		tag, value := tokens[iteration], tokens[iteration+1]
		if !quote.MatchString(value) {
			errors = append(errors, ValidationError{
				Code:    TAG_ERROR,
				Message: fmt.Sprintf("Key %v needs to have a quoted key", tag)},
			)
			iteration = len(tokens)
		}
		if _, present := tagKeys[tag]; present {
			errors = append(errors, ValidationError{
				Code:    TAG_ERROR,
				Message: fmt.Sprintf("Key %v has been repeated", tag)},
			)
			iteration = len(tokens)
		}
		tagKeys[tag] = value
		iteration += 2
	}
	if len(errors) > 0 {
		return errors
	}
	return errors
}
