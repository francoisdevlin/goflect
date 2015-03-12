package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
	"reflect"
)

/*
The view is a basic tool to ensure that only well formed records are every passed to a record source
*/
type view struct {
	match    matcher.Matcher
	delegate privateRecordService
}

func (service view) createAll(record interface{}) error {
	val := reflect.ValueOf(record)

	for i := 0; i < val.Len(); i++ {
		element := val.Index(i).Interface()
		ok, err := service.match.Match(element)
		if err != nil {
			return err
		}
		if !ok {
			return RecordError("Could not create record, does not match")
		}
	}

	return service.delegate.createAll(record)
}

func (service view) readAll(query matcher.Matcher, record ...interface{}) (func(record ...interface{}) bool, error) {
	return service.delegate.readAll(matcher.And(service.match, query), record)
}

func (service view) updateAll(record interface{}, match matcher.Matcher) error {
	ok, err := service.match.Match(record)
	if err != nil {
		return err
	}
	if !ok {
		return RecordError("Could not update record, does not match")
	}
	return service.delegate.updateAll(record, match)
}

func (service view) deleteAll(record interface{}, match matcher.Matcher) error {
	ok, err := service.match.Match(record)
	if err != nil {
		return err
	}
	if !ok {
		return RecordError("Could not delete record, does not match")
	}
	return service.delegate.deleteAll(record, matcher.And(service.match, match))
}
