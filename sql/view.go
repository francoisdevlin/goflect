package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
The view is a basic tool to ensure that only well formed records are every passed to a record source
*/
type view struct {
	match    matcher.Matcher
	delegate RecordService
}

func (service view) Insert(record interface{}) error {
	ok, err := service.match.Match(record)
	if err != nil {
		return err
	}
	if !ok {
		return RecordError("Could not insert record, does not match")
	}
	return service.delegate.Insert(record)
}

func (service view) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service view) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	return service.delegate.readAll(record,
		matcher.And(service.match, match))
}

func (service view) Update(record interface{}) error {
	ok, err := service.match.Match(record)
	if err != nil {
		return err
	}
	if !ok {
		return RecordError("Could not update record, does not match")
	}
	return service.delegate.Update(record)
}

func (service view) Delete(record interface{}) error {
	ok, err := service.match.Match(record)
	if err != nil {
		return err
	}
	if !ok {
		return RecordError("Could not delete record, does not match")
	}
	return service.delegate.Delete(record)
}

func (service view) Restrict(match matcher.Matcher) (RecordService, error) {
	newMatch := matcher.And(service.match, match)
	return view{match: newMatch, delegate: service.delegate}, nil
}
