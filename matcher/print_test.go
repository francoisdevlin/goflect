package matcher

import (
	"fmt"
	"testing"
)

func ExampleNewDefaultPrinter_printBasic() {
	//This is a quick example of how the pretty printer works
	printAll := func(matchers ...Matcher) {
		printer := NewDefaultPrinter()
		for _, matcher := range matchers {
			result, _ := printer.Print(matcher)
			fmt.Println(result)
		}
	}

	printAll(
		Any(),
		None(),
		Eq(int(1)),
		Eq(float64(1)),
		Eq(true),
		Eq("Bacon"),
	)
	//Output:
	//true
	//false
	// _ = 1
	// _ = 1
	// _ = true
	// _ = "Bacon"
}

func ExampleNewDefaultPrinter_printStruct() {
	//This demonstrates how the pretter printer will handle structs with complex expressions
	printer := NewDefaultPrinter()

	matcher := NewStructMatcher()
	matcher.AddField("A", Eq(int(1)))
	matcher.AddField("B", Eq(int(2)))
	result, _ := printer.Print(matcher)
	fmt.Println(result)

	matcher = NewStructMatcher()
	matcher.AddField("A", And(Gte(int(1)), Lte(int(10))))
	matcher.AddField("B", Eq(int(2)))
	result, _ = printer.Print(matcher)
	fmt.Println(result)

	matcher = NewStructMatcher()
	matcher.AddField("A", Eq(matcher.Field("B")))
	result, _ = printer.Print(matcher)
	fmt.Println(result)

	//Output:
	//A = 1 AND B = 2
	//A >= 1 AND A <= 10 AND B = 2
	//A = B
}

func TestDefaultPrinterFields(t *testing.T) {
	assertMatch := func(expected string, matcher Matcher) {
		printer := NewDefaultPrinter()
		result, _ := printer.Print(matcher)
		if result != expected {
			t.Errorf("got:%v, want:%v", result, expected)
		}
	}
	assertMatch("true", Any())
	assertMatch("false", None())
	assertMatch("_ = 1", Eq(1))
	assertMatch("_ = \"1\"", Eq("1"))
	assertMatch("_ != 1", Neq(1))
	assertMatch("_ != \"1\"", Neq("1"))
	assertMatch("_ < 1", Lt(1))
	assertMatch("_ < \"1\"", Lt("1"))
	assertMatch("_ <= 1", Lte(1))
	assertMatch("_ <= \"1\"", Lte("1"))
	assertMatch("_ > 1", Gt(1))
	assertMatch("_ > \"1\"", Gt("1"))
	assertMatch("_ >= 1", Gte(1))
	assertMatch("_ >= \"1\"", Gte("1"))
	assertMatch("_ IN [1 2 3]", In([]int{1, 2, 3}))
	assertMatch("_ IN [\"1\" \"2\" \"3\"]", In([]string{"1", "2", "3"}))
	assertMatch("_ NOT IN [1 2 3]", NotIn([]int{1, 2, 3}))
	assertMatch("_ MATCH \"1\"", Match("1"))
	assertMatch("_ NOT MATCH \"1\"", NotMatch("1"))

	//Demonstrate Not, with rewrite rules
	assertMatch("false", Not(Any()))
	assertMatch("true", Not(None()))
	assertMatch("true", Not(Not(Any())))
	assertMatch("false", Not(Not(None())))
	assertMatch("NOT (_ > 5 AND _ < 10)", Not(And(Gt(5), Lt(10))))
	assertMatch("_ > 5 AND _ < 10", Not(Not(And(Gt(5), Lt(10)))))

	assertMatch("true", And(Any(), Any()))
	assertMatch("true", Or(Any(), Any()))
	m := NewStructMatcher()
	m.AddField("A", Eq(1))
	m.AddField("B", Eq(0))
	assertMatch("A = 1 AND B = 0", m)
	m = NewStructMatcher()
	m.AddField("A", Eq(m.Field("B")))
	assertMatch("A = B", m)
}

func TestSqlitePrinterFields(t *testing.T) {
	assertMatch := func(expected string, matcher Matcher) {
		printer := NewSqlitePrinter()
		result, _ := printer.Print(matcher)
		if result != expected {
			t.Errorf("got:%v, want:%v", result, expected)
		}
	}
	assertMatch("1", Any())
	assertMatch("0", None())
	assertMatch("_ = 1", Eq(1))
	assertMatch("_ = '1'", Eq("1"))
	assertMatch("_ != 1", Neq(1))
	assertMatch("_ != '1'", Neq("1"))
	assertMatch("_ < 1", Lt(1))
	assertMatch("_ < '1'", Lt("1"))
	assertMatch("_ <= 1", Lte(1))
	assertMatch("_ <= '1'", Lte("1"))
	assertMatch("_ > 1", Gt(1))
	assertMatch("_ > '1'", Gt("1"))
	assertMatch("_ >= 1", Gte(1))
	assertMatch("_ >= '1'", Gte("1"))
	assertMatch("_ IN (1, 2, 3)", In([]int{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]int64{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]int32{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]int16{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]int8{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]uint{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]uint64{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]uint32{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]uint16{1, 2, 3}))
	assertMatch("_ IN (1, 2, 3)", In([]uint8{1, 2, 3}))
	assertMatch("_ IN (1.0000, 2.0000, 3.0000)", In([]float64{1, 2, 3}))
	assertMatch("_ IN (1.0000, 2.0000, 3.0000)", In([]float32{1, 2, 3}))
	assertMatch("_ IN ('1', '2', '3')", In([]string{"1", "2", "3"}))
	assertMatch("_ NOT IN (1, 2, 3)", NotIn([]int{1, 2, 3}))
	assertMatch("_ MATCH '1'", Match("1"))
	assertMatch("_ NOT MATCH '1'", NotMatch("1"))

	assertMatch("0", Not(Any()))
	assertMatch("1", Not(None()))
	assertMatch("1", Not(Not(Any())))
	assertMatch("0", Not(Not(None())))

	assertMatch("1", And(Any(), Any()))
	assertMatch("1", Or(Any(), Any()))

	m := NewStructMatcher()
	m.AddField("A", Eq(1))
	m.AddField("B", Eq(0))
	assertMatch("A = 1 AND B = 0", m)

	m = NewStructMatcher()
	m.AddField("A", Eq(m.Field("B")))
	assertMatch("A = B", m)
}

/*
This demonstrates a very basic example of using a struct matcher to write an SQL where clause
*/
func ExampleNewSqlitePrinter_printHello() {
	printMatcher := func(matcher Matcher) {
		printer := NewSqlitePrinter()
		result, _ := printer.Print(matcher)
		fmt.Println(result)
	}
	m := NewStructMatcher()
	m.AddField("A", Eq(1))
	m.AddField("B", Eq(0))

	printMatcher(m)
	//Output:
	//A = 1 AND B = 0
}
