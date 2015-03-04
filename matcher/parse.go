package matcher

import (
	"fmt"
	//"sort"
	"strconv"
	//"strings"
	"reflect"
	"regexp"
)

type parseErrors int

/*
This is the complete list of parse error codes
*/
const (
	VALID parseErrors = iota
	INVALID_OPERATION
	UNFINISHED_MESSAGE
	TOKENIZE_ERROR
	UNKNOWN_FIELD
	PROMOTION_ERROR
)

/*
This is the error struct that will be returned
*/
type MatchParseError struct {
	Code    parseErrors
	Message string
}

func (s MatchParseError) Error() string {
	return s.Message
}

type Parser interface {
	Parse(string) (Matcher, error)
}
type parseStruct struct {
	Fields map[string]reflect.Kind
}

func promoteToInterface(kind reflect.Kind, value string) (interface{}, error) {
	var val interface{} = nil
	var err error = nil
	switch kind {
	case reflect.Float64:
		val, err = strconv.ParseFloat(value, 64)
	case reflect.Float32:
		val, err = strconv.ParseFloat(value, 32)
		val = float32(val.(float64))
	case reflect.Bool:
		val, err = strconv.ParseBool(value)
	case reflect.Int:
		val, err = strconv.ParseInt(value, 10, 64)
		val = int(val.(int64))
	case reflect.Int64:
		val, err = strconv.ParseInt(value, 10, 64)
	case reflect.Int32:
		val, err = strconv.ParseInt(value, 10, 32)
		val = int32(val.(int64))
	case reflect.Int16:
		val, err = strconv.ParseInt(value, 10, 16)
		val = int16(val.(int64))
	case reflect.Int8:
		val, err = strconv.ParseInt(value, 10, 8)
		val = int8(val.(int64))
	case reflect.Uint:
		val, err = strconv.ParseUint(value, 10, 64)
		val = uint(val.(uint64))
	case reflect.Uint64:
		val, err = strconv.ParseUint(value, 10, 64)
	case reflect.Uint32:
		val, err = strconv.ParseUint(value, 10, 32)
		val = uint32(val.(uint64))
	case reflect.Uint16:
		val, err = strconv.ParseUint(value, 10, 16)
		val = uint16(val.(uint64))
	case reflect.Uint8:
		val, err = strconv.ParseUint(value, 10, 8)
		val = uint8(val.(uint64))
	default:
		val = value
	}
	return val, err
}

func (service parseStruct) Parse(input string) (Matcher, error) {
	tokens, err := tokenize(input)
	if err != nil {
		fmt.Println(tokens)
		return nil, err
	}
	//output := structMatcher{}
	output := And()
	Lookup := map[string]fieldOps{
		"=":         EQ,
		"!=":        NEQ,
		"<":         LT,
		"<=":        LTE,
		">":         GT,
		">=":        GTE,
		"IN":        IN,
		"NOT IN":    NOT_IN,
		"MATCH":     MATCH,
		"NOT MATCH": NOT_MATCH,
	}
	iteration := 0

	conjoin := And
	field, op, value := "", "", ""
	cleanParse := VALID
	returnF := func(message string) (Matcher, error) {
		return nil, MatchParseError{Code: cleanParse, Message: message}
	}
	for iteration < len(tokens) {
		switch {
		case tokens[iteration] == "AND":
			conjoin = And
			cleanParse = UNFINISHED_MESSAGE
		case tokens[iteration] == "OR":
			conjoin = Or
			cleanParse = UNFINISHED_MESSAGE
		case field == "":
			field = tokens[iteration]
			if _, present := service.Fields[field]; !present && field != "_" {
				cleanParse = UNKNOWN_FIELD
				return returnF("Unknown Field provided: " + field)
			}
			cleanParse = UNFINISHED_MESSAGE
		case op == "":
			op = tokens[iteration]
			cleanParse = UNFINISHED_MESSAGE
		case op == "NOT":
			op += " " + tokens[iteration]
			cleanParse = UNFINISHED_MESSAGE
		case op == "IN" || op == "NOT IN":
			localIteration := iteration + 1
			vals := make([]string, 0)
			for localIteration < len(tokens) && tokens[localIteration] != ")" {
				vals = append(vals, tokens[localIteration])
				localIteration++
			}
			if localIteration >= len(tokens) {
				return returnF("Could not parse contents of IN clause")
			}
			iteration = localIteration
			step := And()
			realOp, _ := Lookup[op]
			if field == "_" {
				step = And(step, fieldMatcher{Op: realOp, Value: value})
			} else {
				temp := NewStructMatcher()
				temp.AddField(field, fieldMatcher{Op: realOp, Value: value})
				step = And(step, temp)
			}
			output = conjoin(output, step)
			field, op, value = "", "", ""
			cleanParse = VALID
		default:
			value = tokens[iteration]
			step := And()
			realOp, present := Lookup[op]
			if !present {
				cleanParse = INVALID_OPERATION
				return returnF("Operation type is not supported: " + op)
			}

			kind := service.Fields[field]
			val, promotionError := promoteToInterface(kind, value)
			valKind, symbolHit := service.Fields[value]
			if promotionError != nil {
				if !symbolHit {
					cleanParse = PROMOTION_ERROR
					return returnF(fmt.Sprintf("Could not promote field %v to kind %v for value '%v'", field, kind, value))
				} else if valKind != kind {
					cleanParse = PROMOTION_ERROR
					return returnF(fmt.Sprintf("Cannot compare fields %v and %v, they are different kinds", field, value))
				}
			}

			if field == "_" {
				step = And(step, fieldMatcher{Op: realOp, Value: val})
			} else {
				temp := NewStructMatcher()
				if symbolHit {
					temp.AddField(field, fieldMatcher{Op: realOp, Value: temp.Field(value)})
				} else {
					temp.AddField(field, fieldMatcher{Op: realOp, Value: val})
				}
				step = And(step, temp)
			}
			output = conjoin(output, step)
			field, op, value = "", "", ""
			cleanParse = VALID
		}
		iteration++
	}
	switch cleanParse {
	case VALID:
		return output, nil
	case UNFINISHED_MESSAGE:
		return returnF("The message has a trailing entry: " + tokens[iteration-1])
	}
	return output, nil
}

func tokenize(message string) ([]string, error) {
	output := make([]string, 0)
	whitespace, _ := regexp.Compile("^[\\s,]+")
	lParen, _ := regexp.Compile("^\\(")
	rParen, _ := regexp.Compile("^\\)")
	symbol, _ := regexp.Compile("^[a-zA-Z_]\\w*")
	number, _ := regexp.Compile("^-?[0-9]+(\\.[0-9]+)?")
	operators, _ := regexp.Compile("^[!=<>]+")
	quote, _ := regexp.Compile("^\"(?:\\\\?.)*?\"")
	for len(message) > 0 {
		switch {
		case whitespace.MatchString(message):
			message = whitespace.ReplaceAllString(message, "")
		case lParen.MatchString(message):
			output = append(output, lParen.FindString(message))
			message = lParen.ReplaceAllString(message, "")
		case rParen.MatchString(message):
			output = append(output, rParen.FindString(message))
			message = rParen.ReplaceAllString(message, "")
		case symbol.MatchString(message):
			output = append(output, symbol.FindString(message))
			message = symbol.ReplaceAllString(message, "")
		case number.MatchString(message):
			output = append(output, number.FindString(message))
			message = number.ReplaceAllString(message, "")
		case quote.MatchString(message):
			output = append(output, quote.FindString(message))
			message = quote.ReplaceAllString(message, "")
		case operators.MatchString(message):
			output = append(output, operators.FindString(message))
			message = operators.ReplaceAllString(message, "")
		default:
			return output, MatchParseError{Code: TOKENIZE_ERROR, Message: message}
		}
	}
	return output, nil
}

func NewParser(context interface{}) (Parser, error) {
	localContext := make(map[string]reflect.Kind)
	switch c := context.(type) {
	case map[string]reflect.Kind:
		localContext = c
	case reflect.Kind:
		localContext["_"] = c
	case map[string]interface{}:
		for name, value := range c {
			typ := reflect.TypeOf(value)
			localContext[name] = typ.Kind()
		}
	default:
		typ := reflect.TypeOf(c)
		switch typ.Kind() {
		case reflect.Bool,
			reflect.String,
			reflect.Float64,
			reflect.Float32,
			reflect.Uint,
			reflect.Uint64,
			reflect.Uint32,
			reflect.Uint16,
			reflect.Uint8,
			reflect.Int,
			reflect.Int64,
			reflect.Int32,
			reflect.Int16,
			reflect.Int8:
			localContext["_"] = typ.Kind()
		default:
			return nil, nil
		}
	}

	return parseStruct{Fields: localContext}, nil
}
