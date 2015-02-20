package goflect

import (
	//"fmt"
	"reflect"
	"regexp"
)

type FieldOps int

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

type Matcher interface {
	Match(record interface{}) (bool, error)
}

type FieldMatcher struct {
	Op         FieldOps
	Value      interface{}
	fieldCache interface{}
}

type InvalidCompare int

func (i InvalidCompare) Error() string {
	return "Invalid comparoson operation"
}

func (field *FieldMatcher) warmCache() {
	if field.fieldCache != nil {
		return
	}
	switch r := field.Value.(type) {
	case []int:
		temp := make(map[int]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []int64:
		temp := make(map[int64]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []int32:
		temp := make(map[int32]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []int16:
		temp := make(map[int16]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []int8:
		temp := make(map[int8]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []uint:
		temp := make(map[uint]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []uint64:
		temp := make(map[uint64]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []uint32:
		temp := make(map[uint32]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []uint16:
		temp := make(map[uint16]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []uint8:
		temp := make(map[uint8]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []float64:
		temp := make(map[float64]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []float32:
		temp := make(map[float32]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []string:
		temp := make(map[string]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	case []bool:
		temp := make(map[bool]int8)
		for _, val := range r {
			temp[val] = 0
		}
		field.fieldCache = temp
	}
}

func (field FieldMatcher) Match(record interface{}) (bool, error) {
	field.warmCache()
	invert := false
	switch field.Op {
	case NEQ, NOT_IN, GT, GTE, NOT_MATCH:
		invert = true
	}
	switch field.Op {
	case NEQ, EQ:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return invert != false, nil
		}
		switch r := record.(type) {
		case int:
			return invert != (r == field.Value.(int)), nil
		case int64:
			return invert != (r == field.Value.(int64)), nil
		case int32:
			return invert != (r == field.Value.(int32)), nil
		case int16:
			return invert != (r == field.Value.(int16)), nil
		case int8:
			return invert != (r == field.Value.(int8)), nil
		case uint:
			return invert != (r == field.Value.(uint)), nil
		case uint64:
			return invert != (r == field.Value.(uint64)), nil
		case uint32:
			return invert != (r == field.Value.(uint32)), nil
		case uint16:
			return invert != (r == field.Value.(uint16)), nil
		case uint8:
			return invert != (r == field.Value.(uint8)), nil
		case float64:
			return invert != (r == field.Value.(float64)), nil
		case float32:
			return invert != (r == field.Value.(float32)), nil
		case string:
			return invert != (r == field.Value.(string)), nil
		case bool:
			return invert != (r == field.Value.(bool)), nil
		}
	case GTE, LT:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return invert != (r < field.Value.(int)), nil
		case int64:
			return invert != (r < field.Value.(int64)), nil
		case int32:
			return invert != (r < field.Value.(int32)), nil
		case int16:
			return invert != (r < field.Value.(int16)), nil
		case int8:
			return invert != (r < field.Value.(int8)), nil
		case uint:
			return invert != (r < field.Value.(uint)), nil
		case uint64:
			return invert != (r < field.Value.(uint64)), nil
		case uint32:
			return invert != (r < field.Value.(uint32)), nil
		case uint16:
			return invert != (r < field.Value.(uint16)), nil
		case uint8:
			return invert != (r < field.Value.(uint8)), nil
		case float64:
			return invert != (r < field.Value.(float64)), nil
		case float32:
			return invert != (r < field.Value.(float32)), nil
		case string:
			return invert != (r < field.Value.(string)), nil
		}
	case GT, LTE:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return invert != (r <= field.Value.(int)), nil
		case int64:
			return invert != (r <= field.Value.(int64)), nil
		case int32:
			return invert != (r <= field.Value.(int32)), nil
		case int16:
			return invert != (r <= field.Value.(int16)), nil
		case int8:
			return invert != (r <= field.Value.(int8)), nil
		case uint:
			return invert != (r <= field.Value.(uint)), nil
		case uint64:
			return invert != (r <= field.Value.(uint64)), nil
		case uint32:
			return invert != (r <= field.Value.(uint32)), nil
		case uint16:
			return invert != (r <= field.Value.(uint16)), nil
		case uint8:
			return invert != (r <= field.Value.(uint8)), nil
		case float64:
			return invert != (r <= field.Value.(float64)), nil
		case float32:
			return invert != (r <= field.Value.(float32)), nil
		case string:
			return invert != (r <= field.Value.(string)), nil
		}
	case IN, NOT_IN:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value).Elem() {
			return invert, nil
		}
		switch r := record.(type) {
		case int:
			_, present := (field.fieldCache.(map[int]int8))[r]
			return invert != present, nil
		case int64:
			_, present := (field.fieldCache.(map[int64]int8))[r]
			return invert != present, nil
		case int32:
			_, present := (field.fieldCache.(map[int32]int8))[r]
			return invert != present, nil
		case int16:
			_, present := (field.fieldCache.(map[int16]int8))[r]
			return invert != present, nil
		case int8:
			_, present := (field.fieldCache.(map[int8]int8))[r]
			return invert != present, nil
		case uint:
			_, present := (field.fieldCache.(map[uint]int8))[r]
			return invert != present, nil
		case uint64:
			_, present := (field.fieldCache.(map[uint64]int8))[r]
			return invert != present, nil
		case uint32:
			_, present := (field.fieldCache.(map[uint32]int8))[r]
			return invert != present, nil
		case uint16:
			_, present := (field.fieldCache.(map[uint16]int8))[r]
			return invert != present, nil
		case uint8:
			_, present := (field.fieldCache.(map[uint8]int8))[r]
			return invert != present, nil
		case float64:
			_, present := (field.fieldCache.(map[float64]int8))[r]
			return invert != present, nil
		case float32:
			_, present := (field.fieldCache.(map[float32]int8))[r]
			return invert != present, nil
		case string:
			_, present := (field.fieldCache.(map[string]int8))[r]
			return invert != present, nil
		case bool:
			_, present := (field.fieldCache.(map[bool]int8))[r]
			return invert != present, nil
		}
	case MATCH, NOT_MATCH:
		if reflect.TypeOf(record) != reflect.TypeOf(field.Value) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case string:
			exp, err := regexp.Compile(field.Value.(string))
			if err != nil {
				return false, InvalidCompare(1)
			}
			return invert != exp.MatchString(r), nil
		}
	}
	return false, InvalidCompare(1)
}
