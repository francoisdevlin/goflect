package goflect

import (
	"testing"
)

func TestAnyMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matchTrue, _, _ := withMatcher(AnyMatch(1))
	matchTrue("1 is Any", 1)
}

func TestNoneMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	_, matchFalse, _ := withMatcher(NoneMatch(1))
	matchFalse("1 is Any", 1)
}
