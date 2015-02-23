package goflect

type AnyMatch int
type NoneMatch int

type InvertMatch struct {
	M Matcher
}

type AndMatch struct {
	Matchers []Matcher
}

type OrMatch struct {
	Matchers []Matcher
}

func (a AnyMatch) Match(record interface{}) (bool, error) {
	return true, nil
}

func (a NoneMatch) Match(record interface{}) (bool, error) {
	return false, nil
}

func (a InvertMatch) Match(record interface{}) (bool, error) {
	result, err := a.M.Match(record)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (a AndMatch) Match(record interface{}) (bool, error) {
	accum := true
	for _, m := range a.Matchers {
		result, err := m.Match(record)
		if err != nil {
			return false, err
		}
		accum = accum && result
		if !accum {
			break
		}
	}
	return accum, nil
}

func (a OrMatch) Match(record interface{}) (bool, error) {
	accum := false
	for _, m := range a.Matchers {
		result, err := m.Match(record)
		if err != nil {
			return false, err
		}
		accum = accum || result
		if accum {
			break
		}
	}
	return accum, nil
}

func Any() Matcher {
	return AnyMatch(1)
}

func None() Matcher {
	return NoneMatch(1)
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

/*
This function will return a matcher that is the logical inverse of the provided matcher.  Sometimes it will wrap the provided matcher with an Inverter, other times it will perform an optimization in order to keep the call tree as small as possible
*/
func Not(matcher Matcher) Matcher {
	switch r := matcher.(type) {
	case InvertMatch:
		return r.M
	case AnyMatch:
		return None()
	case NoneMatch:
		return Any()
	case FieldMatcher:
		switch r.Op {
		case EQ:
			return Neq(r.Value)
		case NEQ:
			return Eq(r.Value)
		case LT:
			return Gte(r.Value)
		case GTE:
			return Lt(r.Value)
		case LTE:
			return Gt(r.Value)
		case GT:
			return Lte(r.Value)
		case IN:
			return NotIn(r.Value)
		case NOT_IN:
			return In(r.Value)
		case MATCH:
			return NotMatch(r.Value.(string))
		case NOT_MATCH:
			return Match(r.Value.(string))
		default:
			return InvertMatch{M: matcher}
		}
	default:
		return InvertMatch{M: matcher}
	}
}

func And(matchers ...Matcher) Matcher {
	return AndMatch{Matchers: matchers}
}

func Or(matchers ...Matcher) Matcher {
	return OrMatch{Matchers: matchers}
}
