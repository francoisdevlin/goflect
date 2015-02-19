package goflect

import (
	//"fmt"
	"reflect"
)

type FieldOps int
type ComposeOps int

const (
	LT FieldOps = iota
	LTE
	GT
	GTE
	EQ
	NEQ
	IN
	NOT_IN
	MATCH
	NOT_MATCH
)

const (
	NOT ComposeOps = iota
	AND
	OR
	XOR
)

type Matcher interface {
	Match(record interface{}) (bool, error)
}

type FieldMatcher struct {
	Op    FieldOps
	Value interface{}
}

type InvalidCompare int

func (i InvalidCompare) Error() string {
	return "Invalid comparoson operation"
}

type StructMatcher struct {
	Fields map[string]FieldMatcher
}

func (field FieldMatcher) Match(record interface{}) (bool, error) {
	switch field.Op {
	case EQ:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, nil
		}
		switch r := record.(type) {
		case int:
			return r == field.Value.(int), nil
		case int64:
			return r == field.Value.(int64), nil
		case int32:
			return r == field.Value.(int32), nil
		case int16:
			return r == field.Value.(int16), nil
		case int8:
			return r == field.Value.(int8), nil
		case uint:
			return r == field.Value.(uint), nil
		case uint64:
			return r == field.Value.(uint64), nil
		case uint32:
			return r == field.Value.(uint32), nil
		case uint16:
			return r == field.Value.(uint16), nil
		case uint8:
			return r == field.Value.(uint8), nil
		case float64:
			return r == field.Value.(float64), nil
		case float32:
			return r == field.Value.(float32), nil
		case string:
			return r == field.Value.(string), nil
		case bool:
			return r == field.Value.(bool), nil
		}
	case NEQ:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return true, nil
		}
		switch r := record.(type) {
		case int:
			return r != field.Value.(int), nil
		case int64:
			return r != field.Value.(int64), nil
		case int32:
			return r != field.Value.(int32), nil
		case int16:
			return r != field.Value.(int16), nil
		case int8:
			return r != field.Value.(int8), nil
		case uint:
			return r != field.Value.(uint), nil
		case uint64:
			return r != field.Value.(uint64), nil
		case uint32:
			return r != field.Value.(uint32), nil
		case uint16:
			return r != field.Value.(uint16), nil
		case uint8:
			return r != field.Value.(uint8), nil
		case float64:
			return r != field.Value.(float64), nil
		case float32:
			return r != field.Value.(float32), nil
		case string:
			return r != field.Value.(string), nil
		case bool:
			return r != field.Value.(bool), nil
		}
	case LT:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return r < field.Value.(int), nil
		case int64:
			return r < field.Value.(int64), nil
		case int32:
			return r < field.Value.(int32), nil
		case int16:
			return r < field.Value.(int16), nil
		case int8:
			return r < field.Value.(int8), nil
		case uint:
			return r < field.Value.(uint), nil
		case uint64:
			return r < field.Value.(uint64), nil
		case uint32:
			return r < field.Value.(uint32), nil
		case uint16:
			return r < field.Value.(uint16), nil
		case uint8:
			return r < field.Value.(uint8), nil
		case float64:
			return r < field.Value.(float64), nil
		case float32:
			return r < field.Value.(float32), nil
		case string:
			return r < field.Value.(string), nil
		}
	case LTE:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return r <= field.Value.(int), nil
		case int64:
			return r <= field.Value.(int64), nil
		case int32:
			return r <= field.Value.(int32), nil
		case int16:
			return r <= field.Value.(int16), nil
		case int8:
			return r <= field.Value.(int8), nil
		case uint:
			return r <= field.Value.(uint), nil
		case uint64:
			return r <= field.Value.(uint64), nil
		case uint32:
			return r <= field.Value.(uint32), nil
		case uint16:
			return r <= field.Value.(uint16), nil
		case uint8:
			return r <= field.Value.(uint8), nil
		case float64:
			return r <= field.Value.(float64), nil
		case float32:
			return r <= field.Value.(float32), nil
		case string:
			return r <= field.Value.(string), nil
		}
	case GT:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return r > field.Value.(int), nil
		case int64:
			return r > field.Value.(int64), nil
		case int32:
			return r > field.Value.(int32), nil
		case int16:
			return r > field.Value.(int16), nil
		case int8:
			return r > field.Value.(int8), nil
		case uint:
			return r > field.Value.(uint), nil
		case uint64:
			return r > field.Value.(uint64), nil
		case uint32:
			return r > field.Value.(uint32), nil
		case uint16:
			return r > field.Value.(uint16), nil
		case uint8:
			return r > field.Value.(uint8), nil
		case float64:
			return r > field.Value.(float64), nil
		case float32:
			return r > field.Value.(float32), nil
		case string:
			return r > field.Value.(string), nil
		}
	case GTE:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return r >= field.Value.(int), nil
		case int64:
			return r >= field.Value.(int64), nil
		case int32:
			return r >= field.Value.(int32), nil
		case int16:
			return r >= field.Value.(int16), nil
		case int8:
			return r >= field.Value.(int8), nil
		case uint:
			return r >= field.Value.(uint), nil
		case uint64:
			return r >= field.Value.(uint64), nil
		case uint32:
			return r >= field.Value.(uint32), nil
		case uint16:
			return r >= field.Value.(uint16), nil
		case uint8:
			return r >= field.Value.(uint8), nil
		case float64:
			return r >= field.Value.(float64), nil
		case float32:
			return r >= field.Value.(float32), nil
		case string:
			return r >= field.Value.(string), nil
		}
	}
	return false, InvalidCompare(1)
}
