package matcher

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

/*
This shows some basic usage of the NOT compound matcher
*/
func ExampleNot_1() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Not(Eq(1)), 1)
	show(Not(Eq(1)), 2)
	show(Not(And(Gt(1), Lte(3))), 2)
	show(Not(And(Gt(1), Lte(3))), 4)
	show(Not(And(Gt(1), Lte(3))), "Bacon")
	//Output:
	//1 != 1 : false
	//2 != 1 : true
	//NOT (2 > 1 AND 2 <= 3) : false
	//NOT (4 > 1 AND 4 <= 3) : true
	//NOT (Bacon > 1 AND Bacon <= 3) : ERROR
}

func ExampleNot_showCompoundExample() {
	printAll := func(matchers ...Matcher) {
		printer := defaultPrinter{}
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
		printer := defaultPrinter{}
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
		printer := defaultPrinter{}
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
		printer := defaultPrinter{}
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

func ExampleEq() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Eq(1), 1)
	show(Eq(2), 1)
	//The types must be exactly the same, or else this is not true
	show(Eq(int(1)), int64(1))
	//Output:
	//1 = 1 : true
	//1 = 2 : false
	//1 = 1 : false
}

func ExampleNeq() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Neq(1), 1)
	show(Neq(2), 1)
	//If the types are not exactly the same they are considered not equal
	show(Neq(int(1)), int64(1))
	//Output:
	//1 != 1 : false
	//1 != 2 : true
	//1 != 1 : true
}

func ExampleLt() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Lt(1), 1)
	show(Lt(2), 1)
	//If the types are not exactly the same it is an error
	show(Lt(int(1)), int64(1))
	//Output:
	//1 < 1 : false
	//1 < 2 : true
	//1 < 1 : ERROR
}

func ExampleLte() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Lte(1), 1)
	show(Lte(2), 1)
	//If the types are not exactly the same it is an error
	show(Lte(int(1)), int64(1))
	//Output:
	//1 <= 1 : true
	//1 <= 2 : true
	//1 <= 1 : ERROR
}

func ExampleGt() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Gt(1), 1)
	show(Gt(2), 1)
	//If the types are not exactly the same it is an error
	show(Gt(int(1)), int64(1))
	//Output:
	//1 > 1 : false
	//1 > 2 : false
	//1 > 1 : ERROR
}

func ExampleGte() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Gte(1), 1)
	show(Gte(2), 1)
	//If the types are not exactly the same it is an error
	show(Gte(int(1)), int64(1))
	//Output:
	//1 >= 1 : true
	//1 >= 2 : false
	//1 >= 1 : ERROR
}

func ExampleNone() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(None(), 1)
	show(None(), 1)
	//If the types are not exactly the same they are considered not equal
	show(None(), int64(1))
	//Output:
	//false : false
	//false : false
	//false : false
}

func ExampleAny() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Any(), 1)
	show(Any(), 1)
	//If the types are not exactly the same they are considered not equal
	show(Any(), int64(1))
	//Output:
	//true : true
	//true : true
	//true : true
}

func ExampleBuggy() {
	show := func(m Matcher, value interface{}) {
		printer := defaultPrinter{v: fmt.Sprint(value)}
		result, err := m.Match(value)
		statement, _ := printer.Print(m)
		if err != nil {
			fmt.Println(statement, ":", "ERROR")
		} else {
			fmt.Println(statement, ":", result)
		}
	}

	show(Buggy(), 1)
	show(Buggy(), 1)
	show(Buggy(), int64(1))

	//Buggy still behaves like a normal matcher, so rewrite rules apply
	show(And(Buggy(), Eq(1)), 1)
	show(And(Buggy(), Any()), 1)
	show(And(Buggy(), None()), 1)
	show(Or(Buggy(), Eq(1)), 1)
	show(Or(Buggy(), Any()), 1)
	show(Or(Buggy(), None()), 1)
	//Output:
	//error : ERROR
	//error : ERROR
	//error : ERROR
	//error AND 1 = 1 : ERROR
	//error : ERROR
	//false : false
	//error OR 1 = 1 : ERROR
	//true : true
	//error : ERROR
}
