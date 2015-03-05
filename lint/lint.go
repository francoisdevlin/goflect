package lint

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
	"regexp"
	"strconv"
)

/*
This is the set of local error codes
*/
type ErrorCode int

const (
	NOMINAL_MISMATCH ErrorCode = iota
	PRIMARY_MISMATCH
	TAG_ERROR
	NOMINAL_MISCOUNT
	PRIMARY_MISCOUNT
	AUTOINC_ERROR
	UNIQUE_ERROR
	BAD_DEFAULT
	VALIDATOR_PARSE_ERROR
)

func (c ErrorCode) String() string {
	switch c {
	case NOMINAL_MISMATCH:
		return "NOMINAL_MISMATCH"
	case PRIMARY_MISMATCH:
		return "PRIMARY_MISMATCH"
	case TAG_ERROR:
		return "TAG_PARSE_ERROR"
	case NOMINAL_MISCOUNT:
		return "NOMINAL_MISCOUNT"
	case PRIMARY_MISCOUNT:
		return "PRIMARY_MISCOUNT"
	case AUTOINC_ERROR:
		return "AUTOINC_ERROR"
	case UNIQUE_ERROR:
		return "UNIQUE_ERROR"
	case BAD_DEFAULT:
		return "BAD_DEFAULT_VALUE"
	case VALIDATOR_PARSE_ERROR:
		return "VALIDATOR_PARSE_ERROR"
	}
	return ""
}

/*
This is the local error type
*/
type ValidationError struct {
	Code    ErrorCode
	Message string
}

func (v ValidationError) Error() string {
	return v.Message
}

/*
This is the function that is called by the linter binary, and delegate the work approrpiately.  Please read the examples to see the constraints on specific items
*/
func ValidateType(record interface{}, list StructList) []Result {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fieldChecks := []func(f goflect.Info) []error{
		nominal,
		uniqueType,
		autoinc,
		checkDefault,
	}

	recordChecks := []func(f []goflect.Info) []error{
		nominalOnce,
		nominalRequiresPrimary,
		primaryOnce,
		validatorStringParseable,
	}

	output := make([]Result, 0)
	fields := goflect.GetInfo(record)

	for _, recordCheck := range recordChecks {
		errors := recordCheck(fields)
		for _, err := range errors {
			cast, _ := err.(ValidationError)
			cast.Message += " on type " + strconv.Quote(typ.Name())
			output = append(output, Result{Error: cast, Position: list.Position})
		}
	}
	for _, field := range fields {
		pos := list.Position
		if fieldStruct, present := list.Structs[field.Name]; present {
			pos = fieldStruct.Position
		}
		valType, _ := typ.FieldByName(field.Name)
		errors := structTag(string(valType.Tag))
		for _, err := range errors {
			cast, _ := err.(ValidationError)
			cast.Message += " with field " + strconv.Quote(field.Name)
			output = append(output, Result{Error: cast, Position: pos})
		}
		for _, fieldCheck := range fieldChecks {
			errors := fieldCheck(field)
			for _, err := range errors {
				cast, _ := err.(ValidationError)
				cast.Message += " with field " + strconv.Quote(field.Name)
				output = append(output, Result{Error: cast, Position: pos})
			}
		}
	}
	return output
}

/*
This test ensures that there is only one field marked nominal for a given type
*/
func nominalOnce(fields []goflect.Info) []error {
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
This test ensures that there is only one field marked nominal for a given type
*/
func nominalRequiresPrimary(fields []goflect.Info) []error {
	nominalFound, primaryFound := false, false
	for _, field := range fields {
		nominalFound = nominalFound || field.IsNominal
		primaryFound = primaryFound || field.IsPrimary
	}

	errors := make([]error, 0)
	if nominalFound && !primaryFound {
		errors = append(errors, ValidationError{
			Code:    NOMINAL_MISMATCH,
			Message: fmt.Sprintf("There is a nominal field without a primary key"),
		})
	}
	return errors
}

/*
This test ensures that there is only one field marked primary for a given type
*/
func primaryOnce(fields []goflect.Info) []error {
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
This checks to ensure that the validator strings are in fact parseable
*/
func validatorStringParseable(fields []goflect.Info) []error {
	errors := make([]error, 0)
	kinds := make(map[string]reflect.Kind)
	for _, field := range fields {
		kinds[field.Name] = field.Kind
	}
	parser, _ := matcher.NewParser(kinds)
	for _, field := range fields {
		_, err := parser.Parse(field.ValidExpr)
		if err != nil {
			errors = append(errors, ValidationError{
				Code:    VALIDATOR_PARSE_ERROR,
				Message: err.Error(),
			})
		}
	}
	return errors
}

/*
This checks the conditions around a nominal field.  The requirements are that the nominal field is type string, and that the nominal field is unique
*/
func nominal(f goflect.Info) []error {
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
This checks the conditions around a unique field.  The requirements are that a unique field must be either an integer or string type.
*/
func uniqueType(f goflect.Info) []error {
	errors := make([]error, 0)
	if !f.IsUnique {
		return errors
	}
	switch f.Kind {
	case reflect.String,
		reflect.Int,
		reflect.Int64,
		reflect.Int32,
		reflect.Int16,
		reflect.Int8,
		reflect.Uint,
		reflect.Uint64,
		reflect.Uint32,
		reflect.Uint16,
		reflect.Uint8:
		return errors
	}
	errors = append(errors, ValidationError{
		Code:    PRIMARY_MISMATCH,
		Message: fmt.Sprintf("Field %v is marked unique, but is kind %v", f.Name, f.Kind),
	})
	return errors
}

/*
This checks the requirements around an autoincremented field.  It must be int-like and have a primary key
*/
func autoinc(f goflect.Info) []error {
	errors := make([]error, 0)
	if !f.IsAutoincrement {
		return errors
	}
	if !f.IsPrimary {
		errors = append(errors, ValidationError{
			Code:    AUTOINC_ERROR,
			Message: fmt.Sprintf("Marked autoinc, but not primary"),
		})
		return errors
	}
	switch f.Kind {
	case reflect.Int,
		reflect.Int64,
		reflect.Int32,
		reflect.Int16,
		reflect.Int8,
		reflect.Uint,
		reflect.Uint64,
		reflect.Uint32,
		reflect.Uint16,
		reflect.Uint8:
		return errors
	}
	errors = append(errors, ValidationError{
		Code:    AUTOINC_ERROR,
		Message: fmt.Sprintf("Field is marked autoinc, but is kind %v", f.Kind),
	})
	return errors
}

/*
This checks the default value, and ensures that is can be parsed by the appropriate type.  The empty string is always valid.  Anything with kind string is also valid.

Note:  There is a seperate check that will parse any matching constraints, and apply them to the defaults
*/
func checkDefault(f goflect.Info) []error {
	errors := make([]error, 0)
	if f.Default == "" {
		return errors
	}
	var err error = nil
	switch f.Kind {
	case reflect.String:
		return errors
	case reflect.Float64:
		_, err = strconv.ParseFloat(f.Default, 64)
	case reflect.Float32:
		_, err = strconv.ParseFloat(f.Default, 32)
	case reflect.Bool:
		_, err = strconv.ParseBool(f.Default)
	case reflect.Int:
		_, err = strconv.ParseInt(f.Default, 10, 64)
	case reflect.Int64:
		_, err = strconv.ParseInt(f.Default, 10, 64)
	case reflect.Int32:
		_, err = strconv.ParseInt(f.Default, 10, 32)
	case reflect.Int16:
		_, err = strconv.ParseInt(f.Default, 10, 16)
	case reflect.Int8:
		_, err = strconv.ParseInt(f.Default, 10, 8)
	case reflect.Uint:
		_, err = strconv.ParseUint(f.Default, 10, 64)
	case reflect.Uint64:
		_, err = strconv.ParseUint(f.Default, 10, 64)
	case reflect.Uint32:
		_, err = strconv.ParseUint(f.Default, 10, 32)
	case reflect.Uint16:
		_, err = strconv.ParseUint(f.Default, 10, 16)
	case reflect.Uint8:
		_, err = strconv.ParseUint(f.Default, 10, 8)
	default:
		errors = append(errors, ValidationError{
			Code:    BAD_DEFAULT,
			Message: fmt.Sprintf("Unable to determine default logic for kind %v", f.Kind),
		})
	}
	if err != nil {
		errors = append(errors, ValidationError{
			Code:    BAD_DEFAULT,
			Message: fmt.Sprintf("Unable to convert \"%v\" to kind %v", f.Default, f.Kind),
		})
	}
	return errors
}

/*
This checks the syntax of the contents of the struct tag, so that the reflection engine works properly.  It will ensure that each tag is unique, the tag can be parsed, and that the tag is followed by a quoted string
*/
func structTag(message string) []error {
	_, output := parseStructTag(message)
	return output
}

func parseStructTag(message string) (map[string]string, []error) {
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
		return nil, errors
	}
	if len(tokens)%2 != 0 {
		errors = append(errors, ValidationError{Code: TAG_ERROR, Message: "There are the wrong number of tokens present"})
	}
	if len(errors) > 0 {
		return nil, errors
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
		return nil, errors
	}
	fieldLimiter := map[string]func(string) []error{
		goflect.TAG_SQL: flagLimiterFactory(goflect.SQL_FIELDS),
		goflect.TAG_UI:  flagLimiterFactory(goflect.UI_FIELDS),
	}
	for tag, value := range tagKeys {
		if limiter, hit := fieldLimiter[tag]; hit {
			for _, err := range limiter(value) {
				if e, ok := err.(ValidationError); ok {
					e.Message += " for tag " + strconv.Quote(tag)
					err = e
				}
				errors = append(errors, err)
			}
		}
	}
	return tagKeys, errors
}

func flagLimiterFactory(flags []string) func(string) []error {
	orderFlags := func(value string) []error {
		errors := make([]error, 0)
		wrapquotes, _ := regexp.Compile("(^\"|\"$)")
		commas, _ := regexp.Compile(", *")
		value = wrapquotes.ReplaceAllString(value, "")
		entries := commas.Split(value, -1)
		temp := make(map[string]int)
		for _, flag := range flags {
			temp[flag] = 1
		}

		for _, entry := range entries {
			//We have a flag we shouldn't...
			if _, hit := temp[entry]; !hit {
				errors = append(errors, ValidationError{
					Code:    TAG_ERROR,
					Message: fmt.Sprintf("Flag '%v' is not allowed", entry),
				})
			}
		}
		return errors

	}
	return orderFlags
}
