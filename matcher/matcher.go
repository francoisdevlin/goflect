package matcher

import (
	//"fmt"
	"reflect"
	"regexp"
)

type fieldMatcher struct {
	Op         fieldOps
	Value      interface{}
	fieldCache interface{}
}

type InvalidCompare int

func (i InvalidCompare) Error() string {
	return "Invalid comparison operation"
}

func (field *fieldMatcher) warmCache() {
	if field.fieldCache != nil {
		return
	}
	v := field.Value
	if y, ok := v.(Yielder); ok {
		v, _ = y.Yield()
	}
	if v == nil {
		return
	}
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	val := reflect.ValueOf(v)
	if !(val.Kind() == reflect.Array || val.Kind() == reflect.Slice) {
		return
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	temp := reflect.MakeMap(reflect.MapOf(typ.Elem(), reflect.TypeOf(int8(0))))
	zero := reflect.ValueOf(int8(0))
	for i := 0; i < val.Len(); i++ {
		temp.SetMapIndex(val.Index(i), zero)
	}
	field.fieldCache = temp.Interface()
}

func (field fieldMatcher) Match(record interface{}) (bool, error) {
	field.warmCache()
	invert := false
	switch field.Op {
	case NEQ, NOT_IN, GT, GTE, NOT_MATCH:
		invert = true
	}
	v := field.Value
	if y, ok := v.(Yielder); ok {
		v, _ = y.Yield()
	}
	switch field.Op {
	case NEQ, EQ:
		if reflect.TypeOf(record) != reflect.TypeOf(v) {
			return invert != false, nil
		}
		switch r := record.(type) {
		case int:
			return invert != (r == v.(int)), nil
		case int64:
			return invert != (r == v.(int64)), nil
		case int32:
			return invert != (r == v.(int32)), nil
		case int16:
			return invert != (r == v.(int16)), nil
		case int8:
			return invert != (r == v.(int8)), nil
		case uint:
			return invert != (r == v.(uint)), nil
		case uint64:
			return invert != (r == v.(uint64)), nil
		case uint32:
			return invert != (r == v.(uint32)), nil
		case uint16:
			return invert != (r == v.(uint16)), nil
		case uint8:
			return invert != (r == v.(uint8)), nil
		case float64:
			return invert != (r == v.(float64)), nil
		case float32:
			return invert != (r == v.(float32)), nil
		case string:
			return invert != (r == v.(string)), nil
		case bool:
			return invert != (r == v.(bool)), nil
		}
	case GTE, LT:
		if reflect.TypeOf(record) != reflect.TypeOf(v) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return invert != (r < v.(int)), nil
		case int64:
			return invert != (r < v.(int64)), nil
		case int32:
			return invert != (r < v.(int32)), nil
		case int16:
			return invert != (r < v.(int16)), nil
		case int8:
			return invert != (r < v.(int8)), nil
		case uint:
			return invert != (r < v.(uint)), nil
		case uint64:
			return invert != (r < v.(uint64)), nil
		case uint32:
			return invert != (r < v.(uint32)), nil
		case uint16:
			return invert != (r < v.(uint16)), nil
		case uint8:
			return invert != (r < v.(uint8)), nil
		case float64:
			return invert != (r < v.(float64)), nil
		case float32:
			return invert != (r < v.(float32)), nil
		case string:
			return invert != (r < v.(string)), nil
		}
	case GT, LTE:
		if reflect.TypeOf(record) != reflect.TypeOf(v) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case int:
			return invert != (r <= v.(int)), nil
		case int64:
			return invert != (r <= v.(int64)), nil
		case int32:
			return invert != (r <= v.(int32)), nil
		case int16:
			return invert != (r <= v.(int16)), nil
		case int8:
			return invert != (r <= v.(int8)), nil
		case uint:
			return invert != (r <= v.(uint)), nil
		case uint64:
			return invert != (r <= v.(uint64)), nil
		case uint32:
			return invert != (r <= v.(uint32)), nil
		case uint16:
			return invert != (r <= v.(uint16)), nil
		case uint8:
			return invert != (r <= v.(uint8)), nil
		case float64:
			return invert != (r <= v.(float64)), nil
		case float32:
			return invert != (r <= v.(float32)), nil
		case string:
			return invert != (r <= v.(string)), nil
		}
	case IN, NOT_IN:
		if reflect.TypeOf(record) != reflect.TypeOf(v).Elem() {
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
		if reflect.TypeOf(record) != reflect.TypeOf(v) {
			return false, InvalidCompare(1)
		}
		switch r := record.(type) {
		case string:
			exp, err := regexp.Compile(v.(string))
			if err != nil {
				return false, InvalidCompare(1)
			}
			return invert != exp.MatchString(r), nil
		}
	}
	return false, InvalidCompare(1)
}
