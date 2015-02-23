package goflect

import (
	"fmt"
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

func ExampleNot_showCompoundExample() {
	printAll := func(matchers ...Matcher) {
		printer := DefaultPrinter{}
		for _, matcher := range matchers {
			result, _ := printer.Print(matcher)
			fmt.Println(result)
		}
	}

	printAll(
		Not(And(Lt(10), Gt(5))),
		Not(Or(Gt(10), Lt(5))),
	)
	//Output:
	//NOT (_ < 10 AND _ > 5)
	//NOT (_ > 10 OR _ < 5)

}

func ExampleNot_showRewriteRules() {
	printAll := func(matchers ...Matcher) {
		printer := DefaultPrinter{}
		for _, matcher := range matchers {
			result, _ := printer.Print(matcher)
			fmt.Println(result)
		}
	}

	printAll(
		//Invert the constants
		Not(Any()),
		Not(None()),
		//Invert basic comparison
		Not(Eq(1)),
		Not(Neq(1)),
		Not(Lt(1)),
		Not(Lte(1)),
		Not(Gt(1)),
		Not(Gte(1)),
		//Invert the advanced matchers
		Not(In([]int{1, 2, 3})),
		Not(NotIn([]int{1, 2, 3})),
		Not(Match("\\.")),
		Not(NotMatch("\\.")),
		//Remove an inverter
		Not(Not(And(Lt(10), Gt(5)))),
	)

	//Output:
	//false
	//true
	//_ != 1
	//_ = 1
	//_ >= 1
	//_ > 1
	//_ <= 1
	//_ < 1
	//_ NOT IN [1 2 3]
	//_ IN [1 2 3]
	//_ NOT MATCH "\."
	//_ MATCH "\."
	//_ < 10 AND _ > 5
}

func ExampleAnd_showRewriteRules() {
	printAll := func(matchers ...Matcher) {
		printer := DefaultPrinter{}
		for _, matcher := range matchers {
			result, _ := printer.Print(matcher)
			fmt.Println(result)
		}
	}

	printAll(
		//Optimize calls to Any out
		And(Eq(1), Any()),
		//No args returns Any
		And(),
		//One arg returns arg
		And(Eq(1)),
		//None short circuits
		And(Eq(1), None()),
		//And statements consolidated
		And(And(Neq(1), Neq(2)), Neq(3)),
	)

	//Output:
	//_ = 1
	//true
	//_ = 1
	//false
	//_ != 1 AND _ != 2 AND _ != 3
}

func ExampleOr_showRewriteRules() {
	printAll := func(matchers ...Matcher) {
		printer := DefaultPrinter{}
		for _, matcher := range matchers {
			result, _ := printer.Print(matcher)
			fmt.Println(result)
		}
	}

	printAll(
		//Optimize calls to None out
		Or(Eq(1), None()),
		//No args returns None
		Or(),
		//One arg returns arg
		Or(Eq(1)),
		//Any short circuits
		Or(Eq(1), Any()),
		//Or statements consolidated
		Or(Or(Neq(1), Neq(2)), Neq(3)),
	)

	//Output:
	//_ = 1
	//false
	//_ = 1
	//true
	//_ != 1 OR _ != 2 OR _ != 3
}
