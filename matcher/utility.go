package matcher

type anyMatch int
type noneMatch int
type errorMatch int

type invertMatch struct {
	M Matcher
}

type andMatch struct {
	Matchers []Matcher
}

type orMatch struct {
	Matchers []Matcher
}

type lambdaYield func() (interface{}, error)

func (y lambdaYield) Yield() (interface{}, error) {
	return y()
}

/*
This function is designed to convert a lambda to a Yieler
*/
func NewLambdaYield(f func() (interface{}, error)) Yielder {
	return lambdaYield(f)
}

func (a anyMatch) Match(record interface{}) (bool, error) {
	return true, nil
}

func (a noneMatch) Match(record interface{}) (bool, error) {
	return false, nil
}

func (a errorMatch) Match(record interface{}) (bool, error) {
	return false, InvalidCompare(0)
}

func (a invertMatch) Match(record interface{}) (bool, error) {
	result, err := a.M.Match(record)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (a andMatch) Match(record interface{}) (bool, error) {
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

func (a orMatch) Match(record interface{}) (bool, error) {
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

/*
This returns a matcher that always returns true
*/
func Any() Matcher {
	return anyMatch(1)
}

/*
This returns a matcher that always returns false
*/
func None() Matcher {
	return noneMatch(1)
}

/*
This returns a matcher that always returns an error.  Used for testing
*/
func Buggy() Matcher {
	return errorMatch(1)
}

/*
This returns a matcher thay will test that a tested value is equal to a specified one
*/
func Eq(record interface{}) Matcher {
	return fieldMatcher{Op: EQ, Value: record}
}

/*
This returns a matcher thay will test that a tested value is not equal to a specified one
*/
func Neq(record interface{}) Matcher {
	return fieldMatcher{Op: NEQ, Value: record}
}

/*
This returns a matcher thay will test that a tested value is less than a specified one
*/
func Lt(record interface{}) Matcher {
	return fieldMatcher{Op: LT, Value: record}
}

/*
This returns a matcher thay will test that a tested value is less than or equal to a specified one
*/
func Lte(record interface{}) Matcher {
	return fieldMatcher{Op: LTE, Value: record}
}

/*
This returns a matcher thay will test that a tested value is greater than a specified one
*/
func Gt(record interface{}) Matcher {
	return fieldMatcher{Op: GT, Value: record}
}

/*
This returns a matcher thay will test that a tested value is greater than or equal to a specified one
*/
func Gte(record interface{}) Matcher {
	return fieldMatcher{Op: GTE, Value: record}
}

func In(record interface{}) Matcher {
	return fieldMatcher{Op: IN, Value: record}
}

func NotIn(record interface{}) Matcher {
	return fieldMatcher{Op: NOT_IN, Value: record}
}

func Match(record string) Matcher {
	return fieldMatcher{Op: MATCH, Value: record}
}

func NotMatch(record string) Matcher {
	return fieldMatcher{Op: NOT_MATCH, Value: record}
}

/*
This function will return a matcher that is the logical inverse of the provided matcher.  Sometimes it will wrap the provided matcher with an Inverter, other times it will perform an optimization in order to keep the call tree as small as possible
*/
func Not(matcher Matcher) Matcher {
	switch r := matcher.(type) {
	case invertMatch:
		return r.M
	case anyMatch:
		return None()
	case noneMatch:
		return Any()
	case fieldMatcher:
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
			return invertMatch{M: matcher}
		}
	default:
		return invertMatch{M: matcher}
	}
}

/*
This will take a collection of matchers and return the logical AND of each matcher.  This is an optimizing statement, so some rewriting of the rules may occur.

Calling this with no arguments is the same as calling Any()
*/
func And(matchers ...Matcher) Matcher {
	usedMatchers := make([]Matcher, 0)
	for _, matcher := range matchers {
		switch m := matcher.(type) {
		case anyMatch:
			usedMatchers = usedMatchers
		case andMatch:
			for _, childMatch := range m.Matchers {
				usedMatchers = append(usedMatchers, childMatch)
			}
		case noneMatch:
			return None()
		default:
			usedMatchers = append(usedMatchers, m)
		}
	}
	if len(usedMatchers) == 0 {
		return Any()
	}
	if len(usedMatchers) == 1 {
		return usedMatchers[0]
	}
	return andMatch{Matchers: usedMatchers}
}

/*
This will take a collection of matchers and return the logical OR of each matcher.  This is an optimizing statement, so some rewriting of the rules may occur.

Calling this with no arguments is the same as calling None()
*/
func Or(matchers ...Matcher) Matcher {
	usedMatchers := make([]Matcher, 0)
	for _, matcher := range matchers {
		switch m := matcher.(type) {
		case noneMatch:
			usedMatchers = usedMatchers
		case orMatch:
			for _, childMatch := range m.Matchers {
				usedMatchers = append(usedMatchers, childMatch)
			}
		case anyMatch:
			return Any()
		default:
			usedMatchers = append(usedMatchers, m)
		}
	}
	if len(usedMatchers) == 0 {
		return None()
	}
	if len(usedMatchers) == 1 {
		return usedMatchers[0]
	}
	return orMatch{Matchers: usedMatchers}
}

func NewStructMatcher() StructMatcher {
	return new(structMatcher)
}
