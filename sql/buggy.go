package records

import (
	"git.sevone.com/sdevlin/goflect.git/matcher"
)

/*
This is a dummy service intended for use with testing.  It can be consumed both within and external to this package
*/
type buggyService struct {
}

func (service buggyService) insertAll(record interface{}) error {
	return RecordError("Intentional Insert Error")
}

func (service buggyService) readAll(record interface{}, match matcher.Matcher) (func(record interface{}) bool, error) {
	return nil, RecordError("Intentional Read Error")
}

func (service buggyService) updateAll(record interface{}, match matcher.Matcher) error {
	return RecordError("Intentional Update Error")
}

func (service buggyService) deleteAll(record interface{}, match matcher.Matcher) error {
	return RecordError("Intentional Delete Error")
}
