package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This is a dummy service intended for use with testing.  It can be consumed both within and external to this package
*/
type DummyService struct {
	Inserts int
	Updates int
	Reads   int
	Deletes int
}

func (service *DummyService) Insert(record interface{}) error {
	service.Inserts++
	return nil
}

func (service *DummyService) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service *DummyService) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	service.Reads++
	return nil, nil
}

func (service *DummyService) Update(record interface{}) error {
	service.Updates++
	return nil
}

func (service *DummyService) Delete(record interface{}) error {
	service.Deletes++
	return nil
}

func (service *DummyService) Restrict(match matcher.Matcher) (RecordService, error) {
	return view{match: match, delegate: service}, nil
}
