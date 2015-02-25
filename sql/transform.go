package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
The
*/
type transform struct {
	transformer func(interface{}) (interface{}, error)
	delegate    RecordService
}

func (service transform) Insert(record interface{}) error {
	trans, err := service.transformer(record)
	if err != nil {
		return err
	}
	return service.delegate.Insert(trans)
}

func (service transform) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service transform) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	trans, err := service.transformer(record)
	if err != nil {
		return nil, err
	}
	return service.delegate.readAll(trans, match)
}

func (service transform) Update(record interface{}) error {
	trans, err := service.transformer(record)
	if err != nil {
		return err
	}
	return service.delegate.Update(trans)
}

func (service transform) Delete(record interface{}) error {
	trans, err := service.transformer(record)
	if err != nil {
		return err
	}
	return service.delegate.Delete(trans)
}

func (service transform) Restrict(match matcher.Matcher) (RecordService, error) {
	return view{match: match, delegate: service}, nil
}

/*
This creates a new transform service that will route the request to the appropriate service underneath
*/
func NewTransformService(trans func(interface{}) (interface{}, error), deleg RecordService) RecordService {
	return transform{transformer: trans, delegate: deleg}
}
