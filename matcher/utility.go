package goflect

type AnyMatch int
type NoneMatch int

func (a AnyMatch) Match(record interface{}) (bool, error) {
	return true, nil
}

func (a NoneMatch) Match(record interface{}) (bool, error) {
	return false, nil
}
