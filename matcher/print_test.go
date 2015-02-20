package goflect

import (
	"fmt"
)

func ExampleDefaultPrinter_Print_any() {
	//This pretty prints the any matcher
	printer := DefaultPrinter{}
	result, _ := printer.Print(Any())
	fmt.Println(result)
	//Output: true
}

func ExampleDefaultPrinter_Print_none() {
	//This pretty prints the none matcher
	printer := DefaultPrinter{}
	result, _ := printer.Print(None())
	fmt.Println(result)
	//Output: false
}

func ExampleDefaultPrinter_Print_equalInt() {
	//This pretty prints an int
	printer := DefaultPrinter{}
	result, _ := printer.Print(Eq(int(1)))
	fmt.Println(result)
	//Output: _ = 1
}

func ExampleDefaultPrinter_Print_equalTrue() {
	//This pretty prints true
	printer := DefaultPrinter{}
	result, _ := printer.Print(Eq(true))
	fmt.Println(result)
	//Output: _ = true
}

func ExampleDefaultPrinter_Print_equalFalse() {
	//This pretty prints false
	printer := DefaultPrinter{}
	result, _ := printer.Print(Eq(false))
	fmt.Println(result)
	//Output: _ = false
}

func ExampleDefaultPrinter_Print_equalString() {
	//This pretty prints a string
	printer := DefaultPrinter{}
	result, _ := printer.Print(Eq("Bacon"))
	fmt.Println(result)
	//Output: _ = "Bacon"
}

func ExampleDefaultPrinter_Print_and() {
	//This pretty prints a compound and matcher
	printer := DefaultPrinter{}
	result, _ := printer.Print(And(Eq("Bacon"), Any()))
	fmt.Println(result)
	//Output: (_ = "Bacon") AND (true)
}

func ExampleDefaultPrinter_Print_or() {
	//This pretty prints a compound or matcher
	printer := DefaultPrinter{}
	result, _ := printer.Print(Or(Eq("Bacon"), Any()))
	fmt.Println(result)
	//Output: (_ = "Bacon") OR (true)
}

func ExampleDefaultPrinter_Print_not() {
	//This pretty prints a not matcher
	printer := DefaultPrinter{}
	result, _ := printer.Print(Not(Or(Eq("Bacon"), Any())))
	fmt.Println(result)
	//Output: NOT ((_ = "Bacon") OR (true))
}
func ExampleDefaultPrinter_Print_struct() {
	//This pretty prints a compound and matcher
	printer := DefaultPrinter{}
	matcher := StructMatcher{}
	matcher.Fields = make(map[string]Matcher)
	matcher.Fields["A"] = Eq(int(1))
	matcher.Fields["B"] = Eq(int(2))
	result, _ := printer.Print(matcher)
	fmt.Println(result)
	//Output:
	//(A = 1)
	//AND (B = 2)
}
