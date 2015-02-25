package matcher

import (
	//"fmt"
	//"reflect"
	//"regexp"
	"testing"
)

func TestStructMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matcher := StructMatcher{}
	matcher.AddField("A", Eq(int(1)))
	matcher.AddField("B", Eq(int(2)))
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
