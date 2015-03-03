package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This is a dummy service intended for use with testing.  It can be consumed both within and external to this package
*/
type dummyService struct {
	Creates int
	Updates int
	Reads   int
	Deletes int
}

func (service *dummyService) createAll(record interface{}) error {
	service.Creates++
	return nil
}

func (service *dummyService) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	service.Reads++
	return nil, nil
}

func (service *dummyService) updateAll(record interface{}, match matcher.Matcher) error {
	service.Updates++
	return nil
}

func (service *dummyService) deleteAll(record interface{}, match matcher.Matcher) error {
	service.Deletes++
	return nil
}
