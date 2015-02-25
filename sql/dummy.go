package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This is a dummy service intended for use with testing.  It can be consumed both within and external to this package
*/
type dummyService struct {
	Inserts int
	Updates int
	Reads   int
	Deletes int
}

func (service *dummyService) Insert(record interface{}) error {
	service.Inserts++
	return nil
}

func (service *dummyService) ReadAll(record interface{}) (func(record interface{}) bool, error) {
	return service.readAll(record, matcher.Any())
}

func (service *dummyService) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	service.Reads++
	return nil, nil
}

func (service *dummyService) Update(record interface{}) error {
	service.Updates++
	return nil
}

func (service *dummyService) Delete(record interface{}) error {
	service.Deletes++
	return nil
}

func (service *dummyService) Restrict(match matcher.Matcher) (RecordService, error) {
	return view{match: match, delegate: service}, nil
}

func NewDummyService() RecordService {
	return new(dummyService)
}
