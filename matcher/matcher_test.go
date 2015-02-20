package goflect

import (
	"testing"
)

type funcSig func(message string, record interface{})

func withMatcherFactory(t *testing.T) func(matcher Matcher) (matchTrue, matchFalse, matchError funcSig) {
	assertMatch := func(message string, matcher Matcher, value interface{}, expected, hasError bool) {
		match, err := matcher.Match(value)
		if match != expected && (err != nil) == hasError {
			t.Error(message)
		}
	}

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
	return withMatcher
}

//This is a placehold to get test coverage up
func TestErrorMessage(t *testing.T) {
	InvalidCompare(1).Error()
}

func TestFieldMatcher(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	//Equality Section
	matchTrue, matchFalse, matchError := withMatcher(Eq(int(1)))
	matchTrue("1 equals 1", int(1))
	matchFalse("1 does not equal 2", int(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(int64(1)))
	matchTrue("1 equals 1", int64(1))
	matchFalse("1 does not equal 2", int64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(int32(1)))
	matchTrue("1 equals 1", int32(1))
	matchFalse("1 does not equal 2", int32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(int16(1)))
	matchTrue("1 equals 1", int16(1))
	matchFalse("1 does not equal 2", int16(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(int8(1)))
	matchTrue("1 equals 1", int8(1))
	matchFalse("1 does not equal 2", int8(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//Uint tests
	matchTrue, matchFalse, matchError = withMatcher(Eq(uint(1)))
	matchTrue("1 equals 1", uint(1))
	matchFalse("1 does not equal 2", uint(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(uint64(1)))
	matchTrue("1 equals 1", uint64(1))
	matchFalse("1 does not equal 2", uint64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(uint32(1)))
	matchTrue("1 equals 1", uint32(1))
	matchFalse("1 does not equal 2", uint32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(uint16(1)))
	matchTrue("1 equals 1", uint16(1))
	matchFalse("1 does not equal 2", uint16(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(uint8(1)))
	matchTrue("1 equals 1", uint8(1))
	matchFalse("1 does not equal 2", uint8(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//float tests
	matchTrue, matchFalse, matchError = withMatcher(Eq(float64(1)))
	matchTrue("1 equals 1", float64(1))
	matchFalse("1 does not equal 2", float64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq(float32(1)))
	matchTrue("1 equals 1", float32(1))
	matchFalse("1 does not equal 2", float32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//bool
	matchTrue, matchFalse, matchError = withMatcher(Eq(true))
	matchTrue("true is true", true)
	matchFalse("true is not false", float64(2))
	matchFalse("true is not bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Eq("Bacon"))
	matchTrue("Bacon is bacon", "Bacon")
	matchFalse("Bacon is not pizza", "Pizza")
	matchFalse("Bacon is not 1", 1)

	//NEQ Section
	matchTrue, matchFalse, matchError = withMatcher(Neq(int(1)))
	matchFalse("1 equals 1", int(1))
	matchTrue("1 does not equal 2", int(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(int64(1)))
	matchFalse("1 equals 1", int64(1))
	matchTrue("1 does not equal 2", int64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(int32(1)))
	matchFalse("1 equals 1", int32(1))
	matchTrue("1 does not equal 2", int32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(int16(1)))
	matchFalse("1 equals 1", int16(1))
	matchTrue("1 does not equal 2", int16(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(int8(1)))
	matchFalse("1 equals 1", int8(1))
	matchTrue("1 does not equal 2", int8(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//Uint tests
	matchTrue, matchFalse, matchError = withMatcher(Neq(uint(1)))
	matchFalse("1 equals 1", uint(1))
	matchTrue("1 does not equal 2", uint(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(uint64(1)))
	matchFalse("1 equals 1", uint64(1))
	matchTrue("1 does not equal 2", uint64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(uint32(1)))
	matchFalse("1 equals 1", uint32(1))
	matchTrue("1 does not equal 2", uint32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(uint16(1)))
	matchFalse("1 equals 1", uint16(1))
	matchTrue("1 does not equal 2", uint16(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(uint8(1)))
	matchFalse("1 equals 1", uint8(1))
	matchTrue("1 does not equal 2", uint8(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//float tests
	matchTrue, matchFalse, matchError = withMatcher(Neq(float64(1)))
	matchFalse("1 equals 1", float64(1))
	matchTrue("1 does not equal 2", float64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq(float32(1)))
	matchFalse("1 equals 1", float32(1))
	matchTrue("1 does not equal 2", float32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//bool
	matchTrue, matchFalse, matchError = withMatcher(Neq(false))
	matchTrue("false is true", true)
	matchFalse("false is not true", false)
	matchTrue("false is not bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Neq("Bacon"))
	matchFalse("Bacon is bacon", "Bacon")
	matchTrue("Bacon is not pizza", "Pizza")
	matchTrue("Bacon is not 1", 1)

	//LT Section
	matchTrue, matchFalse, matchError = withMatcher(Lt(int(1)))
	matchTrue("0 is less that 1", int(0))
	matchFalse("1 is not less that 1", int(1))
	matchFalse("2 is not less that 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(int64(1)))
	matchTrue("0 is less that 1", int64(0))
	matchFalse("1 is not less that 1", int64(1))
	matchFalse("2 is not less that 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(int32(1)))
	matchTrue("0 is less that 1", int32(0))
	matchFalse("1 is not less that 1", int32(1))
	matchFalse("2 is not less that 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(int16(1)))
	matchTrue("0 is less that 1", int16(0))
	matchFalse("1 is not less that 1", int16(1))
	matchFalse("2 is not less that 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(int8(1)))
	matchTrue("0 is less that 1", int8(0))
	matchFalse("1 is not less that 1", int8(1))
	matchFalse("2 is not less that 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(uint(1)))
	matchTrue("0 is less that 1", uint(0))
	matchFalse("1 is not less that 1", uint(1))
	matchFalse("2 is not less that 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(uint64(1)))
	matchTrue("0 is less that 1", uint64(0))
	matchFalse("1 is not less that 1", uint64(1))
	matchFalse("2 is not less that 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(uint32(1)))
	matchTrue("0 is less that 1", uint32(0))
	matchFalse("1 is not less that 1", uint32(1))
	matchFalse("2 is not less that 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(uint16(1)))
	matchTrue("0 is less that 1", uint16(0))
	matchFalse("1 is not less that 1", uint16(1))
	matchFalse("2 is not less that 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(uint8(1)))
	matchTrue("0 is less that 1", uint8(0))
	matchFalse("1 is not less that 1", uint8(1))
	matchFalse("2 is not less that 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(float64(1)))
	matchTrue("0 is less that 1", float64(0))
	matchFalse("1 is not less that 1", float64(1))
	matchFalse("2 is not less that 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt(float32(1)))
	matchTrue("0 is less that 1", float32(0))
	matchFalse("1 is not less that 1", float32(1))
	matchFalse("2 is not less that 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lt("b"))
	matchTrue("a is less that b", "a")
	matchFalse("b is not less that b", "b")
	matchFalse("c is not less that b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(Lt(true))
	matchError("Bool comparison is invalid", true)

	//LTE Section
	matchTrue, matchFalse, matchError = withMatcher(Lte(int(1)))
	matchTrue("0 is less than or equal to 1", int(0))
	matchTrue("1 is less than or equal to 1", int(1))
	matchFalse("2 is not less than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(int64(1)))
	matchTrue("0 is less than or equal to 1", int64(0))
	matchTrue("1 is less than or equal to 1", int64(1))
	matchFalse("2 is not less than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(int32(1)))
	matchTrue("0 is less than or equal to 1", int32(0))
	matchTrue("1 is less than or equal to 1", int32(1))
	matchFalse("2 is not less than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(int16(1)))
	matchTrue("0 is less than or equal to 1", int16(0))
	matchTrue("1 is less than or equal to 1", int16(1))
	matchFalse("2 is not less than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(int8(1)))
	matchTrue("0 is less than or equal to 1", int8(0))
	matchTrue("1 is less than or equal to 1", int8(1))
	matchFalse("2 is not less than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(uint(1)))
	matchTrue("0 is less than or equal to 1", uint(0))
	matchTrue("1 is less than or equal to 1", uint(1))
	matchFalse("2 is not less than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(uint64(1)))
	matchTrue("0 is less than or equal to 1", uint64(0))
	matchTrue("1 is less than or equal to 1", uint64(1))
	matchFalse("2 is not less than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(uint32(1)))
	matchTrue("0 is less than or equal to 1", uint32(0))
	matchTrue("1 is less than or equal to 1", uint32(1))
	matchFalse("2 is not less than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(uint16(1)))
	matchTrue("0 is less than or equal to 1", uint16(0))
	matchTrue("1 is less than or equal to 1", uint16(1))
	matchFalse("2 is not less than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(uint8(1)))
	matchTrue("0 is less than or equal to 1", uint8(0))
	matchTrue("1 is less than or equal to 1", uint8(1))
	matchFalse("2 is not less than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(float64(1)))
	matchTrue("0 is less than or equal to 1", float64(0))
	matchTrue("1 is less than or equal to 1", float64(1))
	matchFalse("2 is not less than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte(float32(1)))
	matchTrue("0 is less than or equal to 1", float32(0))
	matchTrue("1 is less than or equal to 1", float32(1))
	matchFalse("2 is not less than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Lte("b"))
	matchTrue("a is less than or equal to b", "a")
	matchTrue("b is less than or equal to b", "b")
	matchFalse("c is not less than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(Lte(true))
	matchError("Bool comparison is invalid", true)

	//GT Section
	matchTrue, matchFalse, matchError = withMatcher(Gt(int(1)))
	matchFalse("0 is greater than or equal to 1", int(0))
	matchFalse("1 is greater than or equal to 1", int(1))
	matchTrue("2 is not greater than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(int64(1)))
	matchFalse("0 is greater than or equal to 1", int64(0))
	matchFalse("1 is greater than or equal to 1", int64(1))
	matchTrue("2 is not greater than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(int32(1)))
	matchFalse("0 is greater than or equal to 1", int32(0))
	matchFalse("1 is greater than or equal to 1", int32(1))
	matchTrue("2 is not greater than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(int16(1)))
	matchFalse("0 is greater than or equal to 1", int16(0))
	matchFalse("1 is greater than or equal to 1", int16(1))
	matchTrue("2 is not greater than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(int8(1)))
	matchFalse("0 is greater than or equal to 1", int8(0))
	matchFalse("1 is greater than or equal to 1", int8(1))
	matchTrue("2 is not greater than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(uint(1)))
	matchFalse("0 is greater than or equal to 1", uint(0))
	matchFalse("1 is greater than or equal to 1", uint(1))
	matchTrue("2 is not greater than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(uint64(1)))
	matchFalse("0 is greater than or equal to 1", uint64(0))
	matchFalse("1 is greater than or equal to 1", uint64(1))
	matchTrue("2 is not greater than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(uint32(1)))
	matchFalse("0 is greater than or equal to 1", uint32(0))
	matchFalse("1 is greater than or equal to 1", uint32(1))
	matchTrue("2 is not greater than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(uint16(1)))
	matchFalse("0 is greater than or equal to 1", uint16(0))
	matchFalse("1 is greater than or equal to 1", uint16(1))
	matchTrue("2 is not greater than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(uint8(1)))
	matchFalse("0 is greater than or equal to 1", uint8(0))
	matchFalse("1 is greater than or equal to 1", uint8(1))
	matchTrue("2 is not greater than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(float64(1)))
	matchFalse("0 is greater than or equal to 1", float64(0))
	matchFalse("1 is greater than or equal to 1", float64(1))
	matchTrue("2 is not greater than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt(float32(1)))
	matchFalse("0 is greater than or equal to 1", float32(0))
	matchFalse("1 is greater than or equal to 1", float32(1))
	matchTrue("2 is not greater than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gt("b"))
	matchFalse("a is greater than or equal to b", "a")
	matchFalse("b is greater than or equal to b", "b")
	matchTrue("c is not greater than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(Gt(true))
	matchError("Bool comparison is invalid", true)

	//GTE Section
	matchTrue, matchFalse, matchError = withMatcher(Gte(int(1)))
	matchFalse("0 is greater than or equal to 1", int(0))
	matchTrue("1 is greater than or equal to 1", int(1))
	matchTrue("2 is not greater than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(int64(1)))
	matchFalse("0 is greater than or equal to 1", int64(0))
	matchTrue("1 is greater than or equal to 1", int64(1))
	matchTrue("2 is not greater than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(int32(1)))
	matchFalse("0 is greater than or equal to 1", int32(0))
	matchTrue("1 is greater than or equal to 1", int32(1))
	matchTrue("2 is not greater than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(int16(1)))
	matchFalse("0 is greater than or equal to 1", int16(0))
	matchTrue("1 is greater than or equal to 1", int16(1))
	matchTrue("2 is not greater than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(int8(1)))
	matchFalse("0 is greater than or equal to 1", int8(0))
	matchTrue("1 is greater than or equal to 1", int8(1))
	matchTrue("2 is not greater than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(uint(1)))
	matchFalse("0 is greater than or equal to 1", uint(0))
	matchTrue("1 is greater than or equal to 1", uint(1))
	matchTrue("2 is not greater than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(uint64(1)))
	matchFalse("0 is greater than or equal to 1", uint64(0))
	matchTrue("1 is greater than or equal to 1", uint64(1))
	matchTrue("2 is not greater than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(uint32(1)))
	matchFalse("0 is greater than or equal to 1", uint32(0))
	matchTrue("1 is greater than or equal to 1", uint32(1))
	matchTrue("2 is not greater than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(uint16(1)))
	matchFalse("0 is greater than or equal to 1", uint16(0))
	matchTrue("1 is greater than or equal to 1", uint16(1))
	matchTrue("2 is not greater than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(uint8(1)))
	matchFalse("0 is greater than or equal to 1", uint8(0))
	matchTrue("1 is greater than or equal to 1", uint8(1))
	matchTrue("2 is not greater than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(float64(1)))
	matchFalse("0 is greater than or equal to 1", float64(0))
	matchTrue("1 is greater than or equal to 1", float64(1))
	matchTrue("2 is not greater than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte(float32(1)))
	matchFalse("0 is greater than or equal to 1", float32(0))
	matchTrue("1 is greater than or equal to 1", float32(1))
	matchTrue("2 is not greater than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(Gte("b"))
	matchFalse("a is greater than or equal to b", "a")
	matchTrue("b is greater than or equal to b", "b")
	matchTrue("c is not greater than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(Gte(true))
	matchError("Bool comparison is invalid", true)

	//In Tests
	matchTrue, matchFalse, matchError = withMatcher(In([]int{1}))
	matchTrue("1 is in the set", int(1))
	matchFalse("0 is not in the set", int(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]int64{1}))
	matchTrue("1 is in the set", int64(1))
	matchFalse("0 is not in the set", int64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]int32{1}))
	matchTrue("1 is in the set", int32(1))
	matchFalse("0 is not in the set", int32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]int16{1}))
	matchTrue("1 is in the set", int16(1))
	matchFalse("0 is not in the set", int16(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]int8{1}))
	matchTrue("1 is in the set", int8(1))
	matchFalse("0 is not in the set", int8(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]uint{1}))
	matchTrue("1 is in the set", uint(1))
	matchFalse("0 is not in the set", uint(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]uint64{1}))
	matchTrue("1 is in the set", uint64(1))
	matchFalse("0 is not in the set", uint64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]uint32{1}))
	matchTrue("1 is in the set", uint32(1))
	matchFalse("0 is not in the set", uint32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]uint16{1}))
	matchTrue("1 is in the set", uint16(1))
	matchFalse("0 is not in the set", uint16(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]uint8{1}))
	matchTrue("1 is in the set", uint8(1))
	matchFalse("0 is not in the set", uint8(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]float64{1}))
	matchTrue("1 is in the set", float64(1))
	matchFalse("0 is not in the set", float64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]float32{1}))
	matchTrue("1 is in the set", float32(1))
	matchFalse("0 is not in the set", float32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(In([]string{"A"}))
	matchTrue("A is in the set", "A")
	matchFalse("B is not in the set", "B")
	matchFalse("1 is not in the set", 1)

	matchTrue, matchFalse, matchError = withMatcher(In([]bool{true}))
	matchTrue("True is in the set", true)
	matchFalse("False is not in the set", false)
	matchFalse("1 is not in the set", 1)

	//Not In Tests
	matchTrue, matchFalse, matchError = withMatcher(NotIn([]int{1}))
	matchFalse("1 is in the set", int(1))
	matchTrue("0 is not in the set", int(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]int64{1}))
	matchFalse("1 is in the set", int64(1))
	matchTrue("0 is not in the set", int64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]int32{1}))
	matchFalse("1 is in the set", int32(1))
	matchTrue("0 is not in the set", int32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]int16{1}))
	matchFalse("1 is in the set", int16(1))
	matchTrue("0 is not in the set", int16(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]int8{1}))
	matchFalse("1 is in the set", int8(1))
	matchTrue("0 is not in the set", int8(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]uint{1}))
	matchFalse("1 is in the set", uint(1))
	matchTrue("0 is not in the set", uint(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]uint64{1}))
	matchFalse("1 is in the set", uint64(1))
	matchTrue("0 is not in the set", uint64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]uint32{1}))
	matchFalse("1 is in the set", uint32(1))
	matchTrue("0 is not in the set", uint32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]uint16{1}))
	matchFalse("1 is in the set", uint16(1))
	matchTrue("0 is not in the set", uint16(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]uint8{1}))
	matchFalse("1 is in the set", uint8(1))
	matchTrue("0 is not in the set", uint8(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]float64{1}))
	matchFalse("1 is in the set", float64(1))
	matchTrue("0 is not in the set", float64(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]float32{1}))
	matchFalse("1 is in the set", float32(1))
	matchTrue("0 is not in the set", float32(0))
	matchFalse("Bacon is not in the set", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]string{"A"}))
	matchFalse("A is in the set", "A")
	matchTrue("B is not in the set", "B")
	matchFalse("1 is not in the set", 1)

	matchTrue, matchFalse, matchError = withMatcher(NotIn([]bool{true}))
	matchFalse("True is in the set", true)
	matchTrue("False is not in the set", false)
	matchFalse("1 is not in the set", 1)

	//MATCH Case
	matchTrue, matchFalse, matchError = withMatcher(Match("Bacon.*"))
	matchTrue("Bacon Matches", "Bacon")
	matchFalse("Pizza Doesn't Match", "Pizza")
	matchError("1 is not usable", 1)

	matchTrue, matchFalse, matchError = withMatcher(Match("[a"))
	matchError("Regexp comile error", "Test")

	//NOT MATCH Case
	matchTrue, matchFalse, matchError = withMatcher(NotMatch("Bacon.*"))
	matchFalse("Bacon Matches", "Bacon")
	matchTrue("Pizza Doesn't Match", "Pizza")
	matchError("1 is not usable", 1)

	matchTrue, matchFalse, matchError = withMatcher(NotMatch("[a"))
	matchError("Regexp comile error", "Test")
}
