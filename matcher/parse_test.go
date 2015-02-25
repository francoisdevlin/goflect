package matcher

import (
	"fmt"
	"testing"
)

func TestParseCodes(t *testing.T) {
	render := func(input string, code ParseErrors) {
		p := ParseStruct{Fields: map[string]int{"A": 1, "B": 2, "C": 3}}
		_, e := p.Parse(input)
		if e != nil {
			err, _ := e.(MatchParseError)
			if err.Code != code {
				t.Errorf("Expected Error code %v, got %v", code, err.Code)
			}
		}
	}

	//Basic Operator Checks
	render("", VALID)
	render("_ = 1", VALID)
	render("_ != 1", VALID)
	render("_ < 1", VALID)
	render("_ <= 1", VALID)
	render("_ > 1", VALID)
	render("_ >= 1", VALID)

	render("_ IN (1, 2, 3)", VALID)
	render("_ NOT IN (1, 2, 3)", VALID)

	render("A = 1", VALID)
	render("A = 1 AND B != 2", VALID)

	//The unfinished Messages
	render("A = 1 AND", UNFINISHED_MESSAGE)
	render("A =", UNFINISHED_MESSAGE)
	render("A", UNFINISHED_MESSAGE)
	render("_ IN (1, 2, 3", UNFINISHED_MESSAGE)
	render("_ NOT IN (1, 2, 3", UNFINISHED_MESSAGE)

	//Invalid Operations
	render("_ BACON 1", INVALID_OPERATION)

	//Unknown Fields
	render("D = 1", UNKNOWN_FIELD)
	render("A = D", UNKNOWN_FIELD)

}

func TestTokenize(t *testing.T) {
	expectedLen := func(s string, length int, code ParseErrors) {
		tokens, e := Tokenize(s)
		if e != nil {
			err, _ := e.(MatchParseError)
			if err.Code != code {
				t.Errorf("Expected Error code %v, got %v", code, err.Code)
			}
			return
		}
		tokenCount := len(tokens)
		if tokenCount != length {
			fmt.Println(s)
			fmt.Println(tokens)
			t.Errorf("There are not enough tokens in expr: %v, expected %v, got %v", s, length, tokenCount)
		}
	}
	expectedLen("", 0, VALID)
	expectedLen("Bacon", 1, VALID)

	//Check whitepsace
	expectedLen("\tBacon\t", 1, VALID)
	expectedLen("\tBacon\tPizza", 2, VALID)

	//Commas are whitespace, facilitate tooling
	expectedLen("\tBacon,Pizza", 2, VALID)

	//Check paren tokenizing
	expectedLen("Bacon (Pizza)", 4, VALID)
	expectedLen("Bacon ( Pizza )", 4, VALID)

	//Check operator tokenizing
	expectedLen("Bacon = Pizza", 3, VALID)
	expectedLen("_ = Pizza", 3, VALID)
	expectedLen("_ = 1", 3, VALID)
	expectedLen("Bacon=Pizza", 3, VALID)
	expectedLen("Bacon!=Pizza", 3, VALID)

	//Check number tokenizing
	expectedLen("Bacon != 1", 3, VALID)
	expectedLen("Bacon!=1", 3, VALID)
	expectedLen("Bacon!=-1", 3, VALID)
	expectedLen("Bacon!=1.0", 3, VALID)
	expectedLen("Bacon!=-1.0", 3, VALID)

	//Check quoting
	expectedLen("Bacon!=\"Hi\"", 3, VALID)
	expectedLen("Bacon!=\"Hi\\\"\"", 3, VALID)
	//This has no trailing "
	expectedLen("Bacon!=\"Hi", 3, TOKENIZE_ERROR)
	//This has a trailing \"
	expectedLen("Bacon!=\"Hi\\\"", 3, TOKENIZE_ERROR)
}
