package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
The dispatch is a basic tool to route a request to a record service based on the results of a function
*/
type dispatch struct {
	dispatcher func(interface{}) (int, error)
	delegates  []RecordService
}

func (service dispatch) Insert(record interface{}) error {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return err
	}
	delegate := service.delegates[delegateId]
	return delegate.Insert(record)
}

func (service dispatch) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service dispatch) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return nil, err
	}
	delegate := service.delegates[delegateId]
	return delegate.readAll(record, match)
}

func (service dispatch) Update(record interface{}) error {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return err
	}
	delegate := service.delegates[delegateId]
	return delegate.Update(record)
}

func (service dispatch) Delete(record interface{}) error {
	delegateId, err := service.dispatcher(record)
	if err != nil {
		return err
	}
	delegate := service.delegates[delegateId]
	return delegate.Delete(record)
}

func (service dispatch) Restrict(match matcher.Matcher) (RecordService, error) {
	return view{match: match, delegate: service}, nil
}

/*
This creates a new dispatch service that will route the request to the appropriate service underneath
*/
func NewDispatchService(disp func(interface{}) (int, error), delegs []RecordService) RecordService {
	return dispatch{dispatcher: disp, delegates: delegs}
}
