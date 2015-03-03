package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

/*
The dispatch is a basic tool to route a request to a record service based on the results of a function
*/
type dispatch struct {
	dispatcher func(interface{}) (int, error)
	delegates  []privateRecordService
}

func (service dispatch) createAll(record interface{}) error {
	index := -1
	val := reflect.ValueOf(record)

	for i := 0; i < val.Len(); i++ {
		element := val.Index(i).Interface()
		delegateId, err := service.dispatcher(element)
		switch {
		case err != nil:
			return err
		case index == -1:
			index = delegateId
		case index != delegateId:
			return RecordError("Could not create record, multiple dispatches detected")
		}
		if err != nil {
			return err
		}
	}
	if index == -1 {
		return RecordError("Could not create record, no index detected")
	}

	if index >= len(service.delegates) {
		return RecordError("Dispatch index out of range")
	}
	delegate := service.delegates[index]
	return delegate.createAll(record)
}

func (service dispatch) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return nil, err
	}
	if delegateId >= len(service.delegates) {
		return nil, RecordError("Dispatch index out of range")
	}
	delegate := service.delegates[delegateId]
	return delegate.readAll(record, match)
}

func (service dispatch) updateAll(record interface{}, match matcher.Matcher) error {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return err
	}
	if delegateId >= len(service.delegates) {
		return RecordError("Dispatch index out of range")
	}
	delegate := service.delegates[delegateId]
	return delegate.updateAll(record, match)
}

func (service dispatch) deleteAll(record interface{}, match matcher.Matcher) error {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return err
	}
	if delegateId >= len(service.delegates) {
		return RecordError("Dispatch index out of range")
	}
	delegate := service.delegates[delegateId]
	return delegate.deleteAll(record, match)
}
