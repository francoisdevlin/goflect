package goflect

import (
	"testing"
)

//This is a placehold to get test coverage up
func TestErrorMessage(t *testing.T) {
	InvalidCompare(1).Error()
}

func TestFieldMatcher(t *testing.T) {
	assertMatch := func(message string, matcher FieldMatcher, value interface{}, expected, hasError bool) {
		match, err := matcher.Match(value)
		if match != expected && (err != nil) == hasError {
			t.Error(message)
		}
	}

	type funcSig func(message string, record interface{})
	withMatcher := func(matcher FieldMatcher) (matchTrue, matchFalse, matchError funcSig) {
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

	//Equality Section
	matchTrue, matchFalse, matchError := withMatcher(FieldMatcher{Op: EQ, Value: int(1)})
	matchTrue("1 equals 1", int(1))
	matchFalse("1 does not equal 2", int(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: int64(1)})
	matchTrue("1 equals 1", int64(1))
	matchFalse("1 does not equal 2", int64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: int32(1)})
	matchTrue("1 equals 1", int32(1))
	matchFalse("1 does not equal 2", int32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: int16(1)})
	matchTrue("1 equals 1", int16(1))
	matchFalse("1 does not equal 2", int16(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: int8(1)})
	matchTrue("1 equals 1", int8(1))
	matchFalse("1 does not equal 2", int8(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//Uint tests
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: uint(1)})
	matchTrue("1 equals 1", uint(1))
	matchFalse("1 does not equal 2", uint(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: uint64(1)})
	matchTrue("1 equals 1", uint64(1))
	matchFalse("1 does not equal 2", uint64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: uint32(1)})
	matchTrue("1 equals 1", uint32(1))
	matchFalse("1 does not equal 2", uint32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: uint16(1)})
	matchTrue("1 equals 1", uint16(1))
	matchFalse("1 does not equal 2", uint16(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: uint8(1)})
	matchTrue("1 equals 1", uint8(1))
	matchFalse("1 does not equal 2", uint8(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//float tests
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: float64(1)})
	matchTrue("1 equals 1", float64(1))
	matchFalse("1 does not equal 2", float64(2))
	matchFalse("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: float32(1)})
	matchTrue("1 equals 1", float32(1))
	matchFalse("1 does not equal 2", float32(2))
	matchFalse("1 does not equal bacon", "Bacon")

	//bool
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: true})
	matchTrue("true is true", true)
	matchFalse("true is not false", float64(2))
	matchFalse("true is not bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: EQ, Value: "Bacon"})
	matchTrue("Bacon is bacon", "Bacon")
	matchFalse("Bacon is not pizza", "Pizza")
	matchFalse("Bacon is not 1", 1)

	//NEQ Section
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: int(1)})
	matchFalse("1 equals 1", int(1))
	matchTrue("1 does not equal 2", int(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: int64(1)})
	matchFalse("1 equals 1", int64(1))
	matchTrue("1 does not equal 2", int64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: int32(1)})
	matchFalse("1 equals 1", int32(1))
	matchTrue("1 does not equal 2", int32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: int16(1)})
	matchFalse("1 equals 1", int16(1))
	matchTrue("1 does not equal 2", int16(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: int8(1)})
	matchFalse("1 equals 1", int8(1))
	matchTrue("1 does not equal 2", int8(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//Uint tests
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: uint(1)})
	matchFalse("1 equals 1", uint(1))
	matchTrue("1 does not equal 2", uint(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: uint64(1)})
	matchFalse("1 equals 1", uint64(1))
	matchTrue("1 does not equal 2", uint64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: uint32(1)})
	matchFalse("1 equals 1", uint32(1))
	matchTrue("1 does not equal 2", uint32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: uint16(1)})
	matchFalse("1 equals 1", uint16(1))
	matchTrue("1 does not equal 2", uint16(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: uint8(1)})
	matchFalse("1 equals 1", uint8(1))
	matchTrue("1 does not equal 2", uint8(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//float tests
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: float64(1)})
	matchFalse("1 equals 1", float64(1))
	matchTrue("1 does not equal 2", float64(2))
	matchTrue("1 does not equal bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: float32(1)})
	matchFalse("1 equals 1", float32(1))
	matchTrue("1 does not equal 2", float32(2))
	matchTrue("1 does not equal bacon", "Bacon")

	//bool
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: false})
	matchTrue("false is true", true)
	matchFalse("false is not true", false)
	matchTrue("false is not bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: NEQ, Value: "Bacon"})
	matchFalse("Bacon is bacon", "Bacon")
	matchTrue("Bacon is not pizza", "Pizza")
	matchTrue("Bacon is not 1", 1)

	//LT Section
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: int(1)})
	matchTrue("0 is less that 1", int(0))
	matchFalse("1 is not less that 1", int(1))
	matchFalse("2 is not less that 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: int64(1)})
	matchTrue("0 is less that 1", int64(0))
	matchFalse("1 is not less that 1", int64(1))
	matchFalse("2 is not less that 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: int32(1)})
	matchTrue("0 is less that 1", int32(0))
	matchFalse("1 is not less that 1", int32(1))
	matchFalse("2 is not less that 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: int16(1)})
	matchTrue("0 is less that 1", int16(0))
	matchFalse("1 is not less that 1", int16(1))
	matchFalse("2 is not less that 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: int8(1)})
	matchTrue("0 is less that 1", int8(0))
	matchFalse("1 is not less that 1", int8(1))
	matchFalse("2 is not less that 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: uint(1)})
	matchTrue("0 is less that 1", uint(0))
	matchFalse("1 is not less that 1", uint(1))
	matchFalse("2 is not less that 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: uint64(1)})
	matchTrue("0 is less that 1", uint64(0))
	matchFalse("1 is not less that 1", uint64(1))
	matchFalse("2 is not less that 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: uint32(1)})
	matchTrue("0 is less that 1", uint32(0))
	matchFalse("1 is not less that 1", uint32(1))
	matchFalse("2 is not less that 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: uint16(1)})
	matchTrue("0 is less that 1", uint16(0))
	matchFalse("1 is not less that 1", uint16(1))
	matchFalse("2 is not less that 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: uint8(1)})
	matchTrue("0 is less that 1", uint8(0))
	matchFalse("1 is not less that 1", uint8(1))
	matchFalse("2 is not less that 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: float64(1)})
	matchTrue("0 is less that 1", float64(0))
	matchFalse("1 is not less that 1", float64(1))
	matchFalse("2 is not less that 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: float32(1)})
	matchTrue("0 is less that 1", float32(0))
	matchFalse("1 is not less that 1", float32(1))
	matchFalse("2 is not less that 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: "b"})
	matchTrue("a is less that b", "a")
	matchFalse("b is not less that b", "b")
	matchFalse("c is not less that b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LT, Value: true})
	matchError("Bool comparison is invalid", true)

	//LTE Section
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: int(1)})
	matchTrue("0 is less than or equal to 1", int(0))
	matchTrue("1 is less than or equal to 1", int(1))
	matchFalse("2 is not less than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: int64(1)})
	matchTrue("0 is less than or equal to 1", int64(0))
	matchTrue("1 is less than or equal to 1", int64(1))
	matchFalse("2 is not less than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: int32(1)})
	matchTrue("0 is less than or equal to 1", int32(0))
	matchTrue("1 is less than or equal to 1", int32(1))
	matchFalse("2 is not less than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: int16(1)})
	matchTrue("0 is less than or equal to 1", int16(0))
	matchTrue("1 is less than or equal to 1", int16(1))
	matchFalse("2 is not less than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: int8(1)})
	matchTrue("0 is less than or equal to 1", int8(0))
	matchTrue("1 is less than or equal to 1", int8(1))
	matchFalse("2 is not less than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: uint(1)})
	matchTrue("0 is less than or equal to 1", uint(0))
	matchTrue("1 is less than or equal to 1", uint(1))
	matchFalse("2 is not less than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: uint64(1)})
	matchTrue("0 is less than or equal to 1", uint64(0))
	matchTrue("1 is less than or equal to 1", uint64(1))
	matchFalse("2 is not less than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: uint32(1)})
	matchTrue("0 is less than or equal to 1", uint32(0))
	matchTrue("1 is less than or equal to 1", uint32(1))
	matchFalse("2 is not less than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: uint16(1)})
	matchTrue("0 is less than or equal to 1", uint16(0))
	matchTrue("1 is less than or equal to 1", uint16(1))
	matchFalse("2 is not less than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: uint8(1)})
	matchTrue("0 is less than or equal to 1", uint8(0))
	matchTrue("1 is less than or equal to 1", uint8(1))
	matchFalse("2 is not less than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: float64(1)})
	matchTrue("0 is less than or equal to 1", float64(0))
	matchTrue("1 is less than or equal to 1", float64(1))
	matchFalse("2 is not less than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: float32(1)})
	matchTrue("0 is less than or equal to 1", float32(0))
	matchTrue("1 is less than or equal to 1", float32(1))
	matchFalse("2 is not less than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: "b"})
	matchTrue("a is less than or equal to b", "a")
	matchTrue("b is less than or equal to b", "b")
	matchFalse("c is not less than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: LTE, Value: true})
	matchError("Bool comparison is invalid", true)

	//GT Section
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: int(1)})
	matchFalse("0 is greater than or equal to 1", int(0))
	matchFalse("1 is greater than or equal to 1", int(1))
	matchTrue("2 is not greater than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: int64(1)})
	matchFalse("0 is greater than or equal to 1", int64(0))
	matchFalse("1 is greater than or equal to 1", int64(1))
	matchTrue("2 is not greater than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: int32(1)})
	matchFalse("0 is greater than or equal to 1", int32(0))
	matchFalse("1 is greater than or equal to 1", int32(1))
	matchTrue("2 is not greater than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: int16(1)})
	matchFalse("0 is greater than or equal to 1", int16(0))
	matchFalse("1 is greater than or equal to 1", int16(1))
	matchTrue("2 is not greater than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: int8(1)})
	matchFalse("0 is greater than or equal to 1", int8(0))
	matchFalse("1 is greater than or equal to 1", int8(1))
	matchTrue("2 is not greater than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: uint(1)})
	matchFalse("0 is greater than or equal to 1", uint(0))
	matchFalse("1 is greater than or equal to 1", uint(1))
	matchTrue("2 is not greater than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: uint64(1)})
	matchFalse("0 is greater than or equal to 1", uint64(0))
	matchFalse("1 is greater than or equal to 1", uint64(1))
	matchTrue("2 is not greater than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: uint32(1)})
	matchFalse("0 is greater than or equal to 1", uint32(0))
	matchFalse("1 is greater than or equal to 1", uint32(1))
	matchTrue("2 is not greater than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: uint16(1)})
	matchFalse("0 is greater than or equal to 1", uint16(0))
	matchFalse("1 is greater than or equal to 1", uint16(1))
	matchTrue("2 is not greater than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: uint8(1)})
	matchFalse("0 is greater than or equal to 1", uint8(0))
	matchFalse("1 is greater than or equal to 1", uint8(1))
	matchTrue("2 is not greater than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: float64(1)})
	matchFalse("0 is greater than or equal to 1", float64(0))
	matchFalse("1 is greater than or equal to 1", float64(1))
	matchTrue("2 is not greater than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: float32(1)})
	matchFalse("0 is greater than or equal to 1", float32(0))
	matchFalse("1 is greater than or equal to 1", float32(1))
	matchTrue("2 is not greater than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: "b"})
	matchFalse("a is greater than or equal to b", "a")
	matchFalse("b is greater than or equal to b", "b")
	matchTrue("c is not greater than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GT, Value: true})
	matchError("Bool comparison is invalid", true)

	//GTE Section
	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: int(1)})
	matchFalse("0 is greater than or equal to 1", int(0))
	matchTrue("1 is greater than or equal to 1", int(1))
	matchTrue("2 is not greater than or equal to 1", int(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: int64(1)})
	matchFalse("0 is greater than or equal to 1", int64(0))
	matchTrue("1 is greater than or equal to 1", int64(1))
	matchTrue("2 is not greater than or equal to 1", int64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: int32(1)})
	matchFalse("0 is greater than or equal to 1", int32(0))
	matchTrue("1 is greater than or equal to 1", int32(1))
	matchTrue("2 is not greater than or equal to 1", int32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: int16(1)})
	matchFalse("0 is greater than or equal to 1", int16(0))
	matchTrue("1 is greater than or equal to 1", int16(1))
	matchTrue("2 is not greater than or equal to 1", int16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: int8(1)})
	matchFalse("0 is greater than or equal to 1", int8(0))
	matchTrue("1 is greater than or equal to 1", int8(1))
	matchTrue("2 is not greater than or equal to 1", int8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: uint(1)})
	matchFalse("0 is greater than or equal to 1", uint(0))
	matchTrue("1 is greater than or equal to 1", uint(1))
	matchTrue("2 is not greater than or equal to 1", uint(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: uint64(1)})
	matchFalse("0 is greater than or equal to 1", uint64(0))
	matchTrue("1 is greater than or equal to 1", uint64(1))
	matchTrue("2 is not greater than or equal to 1", uint64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: uint32(1)})
	matchFalse("0 is greater than or equal to 1", uint32(0))
	matchTrue("1 is greater than or equal to 1", uint32(1))
	matchTrue("2 is not greater than or equal to 1", uint32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: uint16(1)})
	matchFalse("0 is greater than or equal to 1", uint16(0))
	matchTrue("1 is greater than or equal to 1", uint16(1))
	matchTrue("2 is not greater than or equal to 1", uint16(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: uint8(1)})
	matchFalse("0 is greater than or equal to 1", uint8(0))
	matchTrue("1 is greater than or equal to 1", uint8(1))
	matchTrue("2 is not greater than or equal to 1", uint8(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: float64(1)})
	matchFalse("0 is greater than or equal to 1", float64(0))
	matchTrue("1 is greater than or equal to 1", float64(1))
	matchTrue("2 is not greater than or equal to 1", float64(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: float32(1)})
	matchFalse("0 is greater than or equal to 1", float32(0))
	matchTrue("1 is greater than or equal to 1", float32(1))
	matchTrue("2 is not greater than or equal to 1", float32(2))
	matchError("1 cannot be compared to Bacon", "Bacon")

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: "b"})
	matchFalse("a is greater than or equal to b", "a")
	matchTrue("b is greater than or equal to b", "b")
	matchTrue("c is not greater than or equal to b", "c")
	matchError("b cannot be compared to 1", 1)

	matchTrue, matchFalse, matchError = withMatcher(FieldMatcher{Op: GTE, Value: true})
	matchError("Bool comparison is invalid", true)
}
