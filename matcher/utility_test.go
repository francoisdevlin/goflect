package goflect

import (
	"testing"
)

func TestAnyMatch(t *testing.T) {
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

	matchTrue, _, _ := withMatcher(AnyMatch(1))
	matchTrue("1 is Any", 1)
}

func TestNoneMatch(t *testing.T) {
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

	_, matchFalse, _ := withMatcher(NoneMatch(1))
	matchFalse("1 is Any", 1)
}
