/*
This is a utility for creating mock records to help with testing
*/
package goflect

import (
	"git.sevone.com/sdevlin/goflect.git/goflect"
	"reflect"
	"strconv"
)

/*
This is used to create mock records, very useful for testing
*/
type RecordMock interface {
	Mock(n int64, record interface{}) interface{}
}

/*
This is an implementation of the RecordMock Interface
*/
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

	fields := goflect.GetInfo(record)
	coerce := func(v interface{}, fVal reflect.Value) {
		localVal := reflect.ValueOf(v)
		if fVal.Type() != localVal.Type() {
			localVal = localVal.Convert(fVal.Type())
		}
		fVal.Set(localVal)
	}
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
			coerce(n != 0, fieldVal)
		case reflect.String:
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
			coerce(temp, fieldVal)
		default:
			coerce(n, fieldVal)
		}
	}
	return record
}
