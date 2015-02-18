package goflect

import (
	//"database/sql"
	//"fmt"
	"reflect"
	"strconv"
	//"strings"
)

type RecordMock interface {
	Mock(n int64, record interface{}) interface{}
}

type MockerStruct struct {
	SkipId        bool
	SkipImmutable bool
}

func (service MockerStruct) Mock(n int64, record interface{}) interface{} {
	typ := reflect.TypeOf(record)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := GetInfo(record)
	for _, field := range fields {
		if field.IsAutoincrement && service.SkipId {
			continue
		}
		if field.IsImmutable && service.SkipImmutable {
			continue
		}
		fieldVal := val.FieldByName(field.Name)
		switch field.Kind {
		case reflect.Bool:
			fieldVal.Set(reflect.ValueOf(n != 0))
		case reflect.Float64:
			fieldVal.Set(reflect.ValueOf(float64(n)))
		case reflect.Float32:
			fieldVal.Set(reflect.ValueOf(float32(n)))
		case reflect.Int:
			fieldVal.Set(reflect.ValueOf(int(n)))
		case reflect.Int64:
			fieldVal.Set(reflect.ValueOf(n))
		case reflect.Int32:
			fieldVal.Set(reflect.ValueOf(int32(n)))
		case reflect.Int16:
			fieldVal.Set(reflect.ValueOf(int16(n)))
		case reflect.Int8:
			fieldVal.Set(reflect.ValueOf(int8(n)))
		case reflect.Uint:
			fieldVal.Set(reflect.ValueOf(uint(n)))
		case reflect.Uint64:
			fieldVal.Set(reflect.ValueOf(uint64(n)))
		case reflect.Uint32:
			fieldVal.Set(reflect.ValueOf(uint32(n)))
		case reflect.Uint16:
			fieldVal.Set(reflect.ValueOf(uint16(n)))
		case reflect.Uint8:
			fieldVal.Set(reflect.ValueOf(uint8(n)))
		default:
			temp := strconv.FormatInt(n, 10)
			switch n {
			case 1:
				temp = temp + "st"
			case 2:
				temp = temp + "nd"
			case 3:
				temp = temp + "rd"
			default:
				temp = temp + "th"
			}
			fieldVal.Set(reflect.ValueOf(temp))
		}
	}
	return record
}
