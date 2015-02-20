package goflect

import (
	//"fmt"
	//"reflect"
	//"regexp"
	"testing"
)

func TestStructMatch(t *testing.T) {
	assertMatch := func(message string, matcher Matcher, value interface{}, expected, hasError bool) {
		match, err := matcher.Match(value)
		if match != expected && (err != nil) == hasError {
			t.Error(message)
		}
	}

	type funcSig func(message string, record interface{})
	withMatcher := func(matcher Matcher) (matchTrue, matchFalse, matchError funcSig) {
		matchTrue = func(message string, record interface{}) {
			assertMatch(message, matcher, record, true, false)
		}
		matchFalse = func(message string, record interface{}) {
			assertMatch(message, matcher, record, false, false)
		}
		matchError = func(message string, record interface{}) {
			assertMatch(message, matcher, record, false, true)
		}
		return matchTrue, matchFalse, matchError
	}

	matcher := StructMatcher{}
	matcher.Fields = make(map[string]Matcher)
	matcher.Fields["A"] = FieldMatcher{Op: EQ, Value: int(1)}
	matcher.Fields["B"] = FieldMatcher{Op: EQ, Value: int(2)}
	matchTrue, matchFalse, matchError := withMatcher(matcher)

	matchTrue("A Well formed map", map[string]interface{}{
		"A": int(1),
		"B": int(2),
	})
	matchFalse("A Well formed invalid map", map[string]interface{}{
		"A": int(1),
		"B": int(1),
	})
	matchError("A poorly formed map", map[string]interface{}{
		"A": int(1),
	})
	matchError("A mismatched type", map[string]interface{}{
		"A": int64(1),
		"B": int(2),
	})

	type WellFormed struct {
		A int
		B int
	}
	matchTrue("A Well formed struct", WellFormed{A: 1, B: 2})
	matchFalse("A Well formed map invalid struct", WellFormed{A: 1, B: 1})

	type Mismatch struct {
		A int
		B int64
	}
	matchError("A mismatched type", Mismatch{A: 1, B: 2})
}
