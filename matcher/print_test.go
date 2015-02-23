package goflect

import (
	"fmt"
	"testing"
)

func ExampleDefaultPrinter_Print_basic() {
	//This is a quick example of how the pretty printer works
	printer := DefaultPrinter{}

	//Print our two constant matchers
	result, _ := printer.Print(Any())
	fmt.Println(result)

	result, _ = printer.Print(None())
	fmt.Println(result)

	//Print a few basic equality matchers
	result, _ = printer.Print(Eq(int(1)))
	fmt.Println(result)

	result, _ = printer.Print(Eq(float64(1)))
	fmt.Println(result)

	result, _ = printer.Print(Eq(true))
	fmt.Println(result)

	result, _ = printer.Print(Eq("Bacon"))
	fmt.Println(result)
	//Output:
	//true
	//false
	// _ = 1
	// _ = 1
	// _ = true
	// _ = "Bacon"
}

func ExampleDefaultPrinter_Print_compound() {
	//This demonstrates how the pretty printer handles expressions compound
	printer := DefaultPrinter{}

	//An AND conditional
	result, _ := printer.Print(And(Eq("Bacon"), Any()))
	fmt.Println(result)

	//An OR conditional
	result, _ = printer.Print(Or(Eq("Bacon"), None()))
	fmt.Println(result)

	result, _ = printer.Print(Not(Or(Eq("Bacon"), Any())))
	fmt.Println(result)

	//A compound expression
	result, _ = printer.Print(Or(And(Eq("Bacon"), Any()), Neq("Pizza")))
	fmt.Println(result)
	//Output:
	//_ = "Bacon" AND true
	//_ = "Bacon" OR false
	//NOT (_ = "Bacon" OR true)
	//(_ = "Bacon" AND true) OR _ != "Pizza"
}

func ExampleDefaultPrinter_Print_struct() {
	//This demonstrates how the pretter printer will handle structs with complex expressions
	printer := DefaultPrinter{}

	matcher := StructMatcher{}
	matcher.AddField("A", Eq(int(1)))
	matcher.AddField("B", Eq(int(2)))
	result, _ := printer.Print(matcher)
	fmt.Println(result)

	matcher = StructMatcher{}
	matcher.AddField("A", And(Gte(int(1)), Lte(int(10))))
	matcher.AddField("B", Eq(int(2)))
	result, _ = printer.Print(matcher)
	fmt.Println(result)

	matcher = StructMatcher{}
	matcher.AddField("A", Eq(Literal("B")))
	result, _ = printer.Print(matcher)
	fmt.Println(result)

	//Output:
	//(A = 1) AND (B = 2)
	//(A >= 1 AND A <= 10) AND (B = 2)
	//(A = B)
}

func TestDefaultPrinterFields(t *testing.T) {
	assertMatch := func(expected string, matcher Matcher) {
		printer := DefaultPrinter{}
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

	assertMatch("true AND true", And(Any(), Any()))
	assertMatch("true OR true", Or(Any(), Any()))
	m := StructMatcher{}
	m.AddField("A", Eq(1))
	m.AddField("B", Eq(0))
	assertMatch("(A = 1) AND (B = 0)", m)
	m = StructMatcher{}
	m.AddField("A", Eq(Literal("B")))
	assertMatch("(A = B)", m)
}

func TestSqlitePrinterFields(t *testing.T) {
	assertMatch := func(expected string, matcher Matcher) {
		printer := SqlitePrinter{}
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
	assertMatch("_ IN (\"1\", \"2\", \"3\")", In([]string{"1", "2", "3"}))
	assertMatch("_ NOT IN (1, 2, 3)", NotIn([]int{1, 2, 3}))
	assertMatch("_ MATCH \"1\"", Match("1"))
	assertMatch("_ NOT MATCH \"1\"", NotMatch("1"))

	assertMatch("false", Not(Any()))
	assertMatch("true", Not(None()))
	assertMatch("true", Not(Not(Any())))
	assertMatch("false", Not(Not(None())))

	assertMatch("true AND true", And(Any(), Any()))
	assertMatch("true OR true", Or(Any(), Any()))
	m := StructMatcher{}
	m.AddField("A", Eq(1))
	m.AddField("B", Eq(0))
	assertMatch("(A = 1) AND (B = 0)", m)
	m = StructMatcher{}
	m.AddField("A", Eq(Literal("B")))
	assertMatch("(A = B)", m)
}
