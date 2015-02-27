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
