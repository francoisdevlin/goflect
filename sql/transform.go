package records

import (
	//"fmt"
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

/*
The transform service transforms a record before sending it to a delegate
*/
type transform struct {
	transformer func(interface{}) (interface{}, error)
	delegate    privateRecordService
}

func (service transform) createAll(record interface{}) error {
	val := reflect.ValueOf(record)

	slice := reflect.MakeSlice(val.Type(), 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		trans, err := service.transformer(val.Index(i).Interface())
		if err != nil {
			return err
		}
		if trans == nil {
			return RecordError("Tranform returned nil")
		}
		slice = reflect.Append(slice, reflect.ValueOf(trans))
	}
	return service.delegate.createAll(slice.Interface())
}

func (service transform) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	trans, err := service.transformer(record)
	if err != nil {
		return nil, err
	}
	if trans == nil {
		return nil, RecordError("Tranform returned nil")
	}
	return service.delegate.readAll(trans, match)
}

func (service transform) updateAll(record interface{}, match matcher.Matcher) error {
	trans, err := service.transformer(record)
	if err != nil {
		return err
	}
	if trans == nil {
		return RecordError("Tranform returned nil")
	}
	return service.delegate.updateAll(trans, match)
}

func (service transform) deleteAll(record interface{}, match matcher.Matcher) error {
	trans, err := service.transformer(record)
	if err != nil {
		return err
	}
	if trans == nil {
		return RecordError("Tranform returned nil")
	}
	return service.delegate.deleteAll(trans, match)
}
