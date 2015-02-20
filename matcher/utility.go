package goflect

type AnyMatch int
type NoneMatch int

func (a AnyMatch) Match(record interface{}) (bool, error) {
	return true, nil
}

func (a NoneMatch) Match(record interface{}) (bool, error) {
	return false, nil
}

func Eq(record interface{}) Matcher {
	return FieldMatcher{Op: EQ, Value: record}
}

func Neq(record interface{}) Matcher {
	return FieldMatcher{Op: NEQ, Value: record}
}

func Lt(record interface{}) Matcher {
	return FieldMatcher{Op: LT, Value: record}
}

func Lte(record interface{}) Matcher {
	return FieldMatcher{Op: LTE, Value: record}
}

func Gt(record interface{}) Matcher {
	return FieldMatcher{Op: GT, Value: record}
}

func Gte(record interface{}) Matcher {
	return FieldMatcher{Op: GTE, Value: record}
}

func In(record interface{}) Matcher {
	return FieldMatcher{Op: IN, Value: record}
}

func NotIn(record interface{}) Matcher {
	return FieldMatcher{Op: NOT_IN, Value: record}
}

func Match(record string) Matcher {
	return FieldMatcher{Op: MATCH, Value: record}
}

func NotMatch(record string) Matcher {
	return FieldMatcher{Op: NOT_MATCH, Value: record}
}
