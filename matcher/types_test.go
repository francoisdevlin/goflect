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
It is possible to parse a matcher from a string.  This is obviously great for user input.  See the parse interface documentation for examples
*/
func ExampleMatcher_2() {
	//We need to give the parser a context.  In this case it is a single field, of type int
	p, _ := NewParser(int(1))

	match, _ := p.Parse("_ = 1")

	result, _ := match.Match(1)
	if result {
		fmt.Println("The parsed matcher is true")
	} else {
		fmt.Println("The parsed matcher is false")
	}
	//Output:
	//The parsed matcher is true
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

/*
The struct matcher can compare values within a record as well.  This is accomplished with the Field method.  It must be called from the matcher that will be doing the comparison, or things will not be bound properly when Match is called
*/
func ExampleStructMatcher_3() {
	type Password struct {
		Current string
		New     string
		Repeat  string
	}

	storedPassword := "old secret"
	matcher := NewStructMatcher()
	matcher.AddField("Current", Eq(storedPassword))
	matcher.AddField("New", And(
		Eq(matcher.Field("Repeat")),
		Neq(matcher.Field("Current")),
	))

	passwordRequest := Password{
		Current: "old secret",
		New:     "new secret",
		Repeat:  "new secret",
	}
	if match, _ := matcher.Match(passwordRequest); match {
		fmt.Println("The password request is well formed")
	}

	passwordRequest = Password{
		Current: "mistake",
		New:     "new secret",
		Repeat:  "new secret",
	}
	if match, _ := matcher.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the current is wrong")
	}

	passwordRequest = Password{
		Current: "old secreet",
		New:     "new secret",
		Repeat:  "wrong password repeated",
	}
	if match, _ := matcher.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the wrong password repeated")
	}

	passwordRequest = Password{
		Current: "old secreet",
		New:     "old secret",
		Repeat:  "old secret",
	}
	if match, _ := matcher.Match(passwordRequest); !match {
		fmt.Println("The password request failed, the secret was not changed")
	}

	//Output:
	//The password request is well formed
	//The password request failed, the current is wrong
	//The password request failed, the wrong password repeated
	//The password request failed, the secret was not changed
}

/*
One of the first things that make matchers different that using an ordinary lambda is that they can be rendered as WHERE clauses for SQL queries.  This can be done by constructing the appropriate SQL printer, and calling its print method.
*/
func ExampleStructMatcher_4() {
	idMatcher := NewStructMatcher()
	idMatcher.AddField("Id", Eq(1))

	valueMatcher := NewStructMatcher()
	valueMatcher.AddField("Value", Eq(1))

	printer := NewSqlitePrinter()
	result, _ := printer.Print(idMatcher)
	fmt.Println(result)

	result, _ = printer.Print(valueMatcher)
	fmt.Println(result)

	//Here you can see that composition of filters is respected, producing the result you would expect
	result, _ = printer.Print(And(idMatcher, valueMatcher))
	fmt.Println(result)

	fieldMatcher := NewStructMatcher()
	fieldMatcher.AddField("A", Eq(fieldMatcher.Field("B")))

	//Field bindings are respected as well
	result, _ = printer.Print(fieldMatcher)
	fmt.Println(result)
	//Output:
	//Id = 1
	//Value = 1
	//Id = 1 AND Value = 1
	//A = B
}

/*
Matchers can also be parsed.  This has obvios advantages for real time user input.  See the parse method for several examples of expressions you can use
*/
func ExampleStructMatcher_5() {
	//We need to give the parser a context.  In this case it three fields, of type int,int and string
	p, _ := NewParser(map[string]interface{}{
		"A":    int(1),
		"B":    int(1),
		"Name": "",
	})

	type Foo struct {
		A    int
		B    int
		Name string
	}

	match, err := p.Parse("A = 1")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := match.Match(Foo{A: 1})
	fmt.Println("The result is:", result)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Output:
	//The result is: true
}

/*
In addition to being able to construct a WHERE clause from a matcher, it is possible to parse a matcher from a string.  This is obviously great for user input
*/
func ExamplMatcher_2() {
	//We need to give the parser a context.  In this case it is a single field, of type int
	p, _ := NewParser(int(1))

	match, _ := p.Parse("_ = 1")

	result, _ := match.Match(1)
	if result {
		fmt.Println("The parsed matcher is true")
	} else {
		fmt.Println("The parsed matcher is false")
	}
	//Output:
	//The parsed matcher is true

}
