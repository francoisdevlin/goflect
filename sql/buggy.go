package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This is a dummy service intended for use with testing.  It can be consumed both within and external to this package
*/
type buggyService struct {
}

func (service buggyService) createAll(record interface{}) error {
	return RecordError("Intentional Create Error")
}

func (service buggyService) readAll(query matcher.Matcher, record ...interface{}) (func(record ...interface{}) bool, error) {
	return nil, RecordError("Intentional Read Error")
}

func (service buggyService) updateAll(record interface{}, match matcher.Matcher) error {
	return RecordError("Intentional Update Error")
}

func (service buggyService) deleteAll(record interface{}, match matcher.Matcher) error {
	return RecordError("Intentional Delete Error")
}
