package goflect

import (
	"testing"
)

func TestAnyMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matchTrue, _, _ := withMatcher(Any())
	matchTrue("1 is Any", 1)
}

func TestNoneMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	_, matchFalse, _ := withMatcher(None())
	matchFalse("1 is not None", 1)
}

func TestInvertMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matchTrue, _, _ := withMatcher(Not(None()))
	matchTrue("1 is Any", 1)
	_, matchFalse, _ := withMatcher(Not(Any()))
	matchFalse("1 is Any", 1)
}

func TestAndMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matchTrue, matchFalse, matchError := withMatcher(And(Gt(int(1)), Lt(int(10))))
	matchFalse("1 is too small", int(1))
	matchTrue("2 is okay", int(2))
	matchTrue("9 is okay", int(9))
	matchFalse("10 is too big", int(10))
	matchError("Bacon is right out", "Bacon")
}

func TestOrMatch(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	matchTrue, matchFalse, matchError := withMatcher(Or(Lte(int(1)), Gte(int(10))))
	matchTrue("1 is okay", int(1))
	matchFalse("2 is bad", int(2))
	matchFalse("9 is bad", int(9))
	matchTrue("10 is okay", int(10))
	matchError("Bacon is right out", "Bacon")
}
