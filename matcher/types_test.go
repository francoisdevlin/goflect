package matcher

import (
	"fmt"
)

/*
This is the most basic example of using a matcher.  We show how using a matcher is equivalent to an if statement
*/
func ExampleMatcher_1() {
	//This is a convential if statement
	x := 5
	if x == 5 {
		fmt.Println("Hello Old World")
	}

	matcher := Eq(5)
	if match, _ := matcher.Match(5); match {
		fmt.Println("Hello New World")
	}
	//Output:
	//Hello Old World
	//Hello New World
}

/*
This example show the most basic usage of a struct matcher.  You can see a comporable if statement for the matcher
*/
func ExampleStructMatcher_1() {
	type Foo struct {
		Id    int
		Value int
		A     string
		B     string
	}

	sample := Foo{Id: 1, Value: 1, A: "A string", B: "Another string"}

	if sample.Id == 1 && sample.Value == 1 {
		fmt.Println("Old school match")
	}

	matcher := NewStructMatcher()
	matcher.AddField("Id", Eq(1))
	matcher.AddField("Value", Eq(1))

	if match, _ := matcher.Match(sample); match {
		fmt.Println("New school match")
	}
	//Output:
	//Old school match
	//New school match
}

/*
Matchers are run time duck typed.  This example shows how you can apply a matcher to two different structs, as well as a map of interfaces.
*/
func ExampleStructMatcher_2() {
	type Foo struct {
		Id    int
		Value int
		A     string
		B     string
	}

	type Bar struct {
		Id    int
		Value int
		Bacon string
	}

	matcher := NewStructMatcher()
	matcher.AddField("Id", Eq(1))
	matcher.AddField("Value", Eq(1))

	foo := Foo{Id: 1, Value: 1}
	if match, _ := matcher.Match(foo); match {
		fmt.Println("Record foo matches")
	}

	bar := Bar{Id: 1, Value: 1}
	if match, _ := matcher.Match(bar); match {
		fmt.Println("Record bar matches")
	}

	dict := map[string]interface{}{
		"Id":    1,
		"Value": 1,
	}
	if match, _ := matcher.Match(dict); match {
		fmt.Println("The dictionary matches")
	}

	//Output:
	//Record foo matches
	//Record bar matches
	//The dictionary matches
}
