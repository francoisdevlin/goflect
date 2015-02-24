package goflect

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git"
	"go/token"
	"reflect"
)

type ErrorCode int

const (
	NOMINAL_MISMATCH ErrorCode = iota
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
	//fieldChecks := make([]func(f goflect.Info) []error, 0)
	fieldChecks := []func(f goflect.Info) []error{
		Nominal,
	}

	recordChecks := make([]func(f []goflect.Info) []error, 0)

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
