package matcher

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseCodes(t *testing.T) {
	render := func(input string, code parseErrors) {
		p := parseStruct{Fields: map[string]reflect.Kind{
			"A": reflect.String,
			"B": reflect.String,
			"C": reflect.Int,
		}}
		_, e := p.Parse(input)
		if e != nil {
			err, _ := e.(MatchParseError)
			if err.Code != code {
				t.Errorf("Expected Error code %v, got %v for:%v", code, err.Code, input)
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
	render("A = B", VALID)
	render("A = A", VALID)
	render("C = C", VALID)
	render("A = \"B\"", VALID)
	render("A = 1 AND B != 2", VALID)

	//And here come the parens...
	render("(A = 1)", VALID)
	render("(A = 2) OR (A = 1)", VALID)
	render("((A = 2)) OR (A = 1)", VALID)
	render("(((A = 2))) OR (A = 1)", VALID)
	render("((((A = 2))) OR (A = 1))", VALID)

	//Not requires a paren
	render("NOT (A = 1)", VALID)
	render("(A = 2) OR NOT (A = 1)", VALID)
	render("NOT ((((A = 2))) OR (A = 1))", VALID)

	//The unfinished Messages
	render("A = 1 AND", UNFINISHED_MESSAGE)
	render("A =", UNFINISHED_MESSAGE)
	render("A", UNFINISHED_MESSAGE)
	render("_ IN (1, 2, 3", UNFINISHED_MESSAGE)
	render("_ NOT IN (1, 2, 3", UNFINISHED_MESSAGE)
	render("( A = 1", UNFINISHED_MESSAGE)
	render("(( A = 1 )", UNFINISHED_MESSAGE)
	render("NOT A = 1", UNFINISHED_MESSAGE)
	render("NOT", UNFINISHED_MESSAGE)

	//Invalid Operations
	render("_ BACON 1", INVALID_OPERATION)

	//Unknown Fields
	render("D = 1", UNKNOWN_FIELD)
	render("A = D", UNKNOWN_FIELD)
	render("A = 1 )", UNKNOWN_FIELD) //Dangling parens show up as unknown fields

	//Promotion Error
	render("A = C", PROMOTION_ERROR)
	render("C = A", PROMOTION_ERROR)
	render("C = \"Fail\"", PROMOTION_ERROR)
}

func TestWhitespacePermutations(t *testing.T) {

	//This determine if the indicies have been maximized
	workTodo := func(size int, entries []int) bool {
		for _, value := range entries {
			if value <= size {
				return true
			} else {
				return false
			}
		}
		return false
	}

	//This takes the next step in the tree
	next := func(size int, entries []int) []int {
		for i := len(entries) - 1; i >= 0; i-- {
			if (entries[i]) >= size {
				if i > 0 {
					entries[i] = 0
					entries[i-1]++
				}
			} else {
				entries[i]++
				return entries
			}
		}
		return entries
	}

	tryTokens := func(inputTokens, inputWhitespace []string) {
		size := len(inputWhitespace) - 1
		indicies := make([]int, len(inputTokens)+1)
		entries := make([]string, 2*len(inputTokens)+1)
		for workTodo(size, indicies) {
			entries[0] = inputWhitespace[indicies[0]]
			for i, token := range inputTokens {
				entries[2*i+1] = token
				entries[2*i+2] = inputWhitespace[indicies[i+1]]
			}
			testString := strings.Join(entries, "")
			newTokens, err := tokenize(testString)
			if err != nil {
				t.Errorf("Unexpected token error (%v) with input %v ", err, testString)
			}
			if len(newTokens) != len(inputTokens) {
				t.Errorf("Token error, got %v, want %v", newTokens, inputTokens)
			}
			if strings.Join(newTokens, " ") != strings.Join(inputTokens, " ") {
				t.Errorf("Token error, got %v, want %v", newTokens, inputTokens)
			}
			indicies = next(size, indicies)
		}
	}
	whitespaceNoEmpty := []string{" ", ",", " ,\t\n"}
	whitespace := append(whitespaceNoEmpty, "")
	tryTokens([]string{"A", "=", "1"}, whitespace)
	tryTokens([]string{"(", "A", "=", "1", ")"}, whitespace)
	tryTokens([]string{"(", "(", "A", "=", "1", ")", ")"}, whitespace)
	tryTokens([]string{"A", "=", "1", "AND", "B", "=", "1"}, whitespaceNoEmpty)
}

func TestParsing(t *testing.T) {
	withMatcher := withMatcherFactory(t)

	determineResults := func(op fieldOps, matcher Matcher) (A, B, C funcSig) {
		T, F, E := withMatcher(matcher)
		switch op {
		case EQ:
			return T, F, F
		case NEQ:
			return F, T, T
		case LT:
			return F, F, E
		case LTE:
			return T, F, E
		case GT:
			return F, T, E
		case GTE:
			return T, T, E
		default:
			return T, F, E
		}
	}

	comparisonWorkout := func(parser Parser, smaller, bigger, nonsense interface{}) {
		matcher, _ := parser.Parse("")
		_, ok := matcher.(anyMatch)
		if !ok {
			t.Error("Expected an ANY matcher")
		}
		ops := []fieldOps{EQ, NEQ, LT, LTE, GT, GTE}
		for _, op := range ops {
			input := fmt.Sprintf("_ %v %v", op, smaller)
			matcher, _ = parser.Parse(input)
			A, B, C := determineResults(op, matcher)
			A(fmt.Sprintf("Matching %v,%v for %v", smaller, smaller, op), smaller)
			B(fmt.Sprintf("Matching %v,%v for %v", smaller, bigger, op), bigger)
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nonsense, op), nonsense)
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nil, op), nil)
		}
	}

	fieldName := "_"
	moreWorkouts := func(k reflect.Kind, target, smaller, bigger, nonsense interface{}) {
		p, _ := NewParser(map[string]reflect.Kind{fieldName: k})
		comparisonWorkout(p, smaller, bigger, nonsense)
		p, _ = NewParser(k)
		comparisonWorkout(p, smaller, bigger, nonsense)

		p, _ = NewParser(map[string]interface{}{fieldName: target})
		comparisonWorkout(p, smaller, bigger, nonsense)
		p, _ = NewParser(target)
		comparisonWorkout(p, smaller, bigger, nonsense)
	}

	moreWorkouts(reflect.Int, int(1), int(1), int(2), "Bacon")
	moreWorkouts(reflect.Int64, int64(1), int64(1), int64(2), "Bacon")
	moreWorkouts(reflect.Int32, int32(1), int32(1), int32(2), "Bacon")
	moreWorkouts(reflect.Int16, int16(1), int16(1), int16(2), "Bacon")
	moreWorkouts(reflect.Int8, int8(1), int8(1), int8(2), "Bacon")

	moreWorkouts(reflect.Uint, uint(1), uint(1), uint(2), "Bacon")
	moreWorkouts(reflect.Uint64, uint64(1), uint64(1), uint64(2), "Bacon")
	moreWorkouts(reflect.Uint32, uint32(1), uint32(1), uint32(2), "Bacon")
	moreWorkouts(reflect.Uint16, uint16(1), uint16(1), uint16(2), "Bacon")
	moreWorkouts(reflect.Uint8, uint8(1), uint8(1), uint8(2), "Bacon")

	moreWorkouts(reflect.Float64, float64(1), float64(1), float64(2), "Bacon")
	moreWorkouts(reflect.Float32, float32(1), float32(1), float32(2), "Bacon")

	buildMap := func(value interface{}) map[string]interface{} {
		return map[string]interface{}{
			"A": value,
		}
	}

	comparisonWorkout = func(parser Parser, smaller, bigger, nonsense interface{}) {
		matcher, _ := parser.Parse("")
		_, ok := matcher.(anyMatch)
		if !ok {
			t.Error("Expected an ANY matcher")
		}
		ops := []fieldOps{EQ, NEQ, LT, LTE, GT, GTE}
		for _, op := range ops {
			input := fmt.Sprintf("A %v %v", op, smaller)
			matcher, err := parser.Parse(input)
			if err != nil {
				t.Error(err.Error())
				return
			}
			A, B, C := determineResults(op, matcher)
			A(fmt.Sprintf("Matching %v,%v for %v", smaller, smaller, op), buildMap(smaller))
			B(fmt.Sprintf("Matching %v,%v for %v", smaller, bigger, op), buildMap(bigger))
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nonsense, op), buildMap(nonsense))
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nil, op), buildMap(nil))
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nil, op), nil)
		}
	}

	fieldName = "A"
	//Note - We can't test the annonymous values here, they are supposed to fail
	moreWorkouts = func(k reflect.Kind, target, smaller, bigger, nonsense interface{}) {
		p, _ := NewParser(map[string]reflect.Kind{fieldName: k})
		comparisonWorkout(p, smaller, bigger, nonsense)

		p, _ = NewParser(map[string]interface{}{fieldName: target})
		comparisonWorkout(p, smaller, bigger, nonsense)
	}
	moreWorkouts(reflect.Int, int(1), int(1), int(2), "Bacon")
	moreWorkouts(reflect.Int64, int64(1), int64(1), int64(2), "Bacon")
	moreWorkouts(reflect.Int32, int32(1), int32(1), int32(2), "Bacon")
	moreWorkouts(reflect.Int16, int16(1), int16(1), int16(2), "Bacon")
	moreWorkouts(reflect.Int8, int8(1), int8(1), int8(2), "Bacon")

	moreWorkouts(reflect.Uint, uint(1), uint(1), uint(2), "Bacon")
	moreWorkouts(reflect.Uint64, uint64(1), uint64(1), uint64(2), "Bacon")
	moreWorkouts(reflect.Uint32, uint32(1), uint32(1), uint32(2), "Bacon")
	moreWorkouts(reflect.Uint16, uint16(1), uint16(1), uint16(2), "Bacon")
	moreWorkouts(reflect.Uint8, uint8(1), uint8(1), uint8(2), "Bacon")

	moreWorkouts(reflect.Float64, float64(1), float64(1), float64(2), "Bacon")
	moreWorkouts(reflect.Float32, float32(1), float32(1), float32(2), "Bacon")

	fieldName = "_"
	comparisonWorkout = func(parser Parser, smaller, bigger, nonsense interface{}) {
		matcher, _ := parser.Parse("")
		_, ok := matcher.(anyMatch)
		if !ok {
			t.Error("Expected an ANY matcher")
		}
		ops := []fieldOps{EQ, NEQ, LT, LTE, GT, GTE}
		for _, op := range ops {
			input := fmt.Sprintf("_ %v \"%v\"", op, smaller)
			matcher, err := parser.Parse(input)
			if err != nil {
				t.Errorf("Unexpected parse error '%v'", input)
			}
			A, B, C := determineResults(op, matcher)
			A(fmt.Sprintf("Matching %v,%v for %v", smaller, smaller, op), smaller)
			B(fmt.Sprintf("Matching %v,%v for %v", smaller, bigger, op), bigger)
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nonsense, op), nonsense)
			C(fmt.Sprintf("Matching %v,%v for %v", smaller, nil, op), nil)
		}
	}

	moreWorkouts = func(k reflect.Kind, target, smaller, bigger, nonsense interface{}) {
		p, _ := NewParser(map[string]reflect.Kind{fieldName: k})
		comparisonWorkout(p, smaller, bigger, nonsense)
		p, _ = NewParser(k)
		comparisonWorkout(p, smaller, bigger, nonsense)

		p, _ = NewParser(map[string]interface{}{fieldName: target})
		comparisonWorkout(p, smaller, bigger, nonsense)
		p, _ = NewParser(target)
		comparisonWorkout(p, smaller, bigger, nonsense)
	}

	moreWorkouts(reflect.String, "A", "A", "B", 1)

	//Strings unknown ATM...
	//p = parseStruct{Fields: map[string]reflect.Kind{
	//fieldName: reflect.String,
	//}}
	//comparisonWorkout(p, "A", "B", 1)

	//p = parseStruct{Fields: map[string]reflect.Kind{
	//fieldName: reflect.String,
	//}}
	//comparisonWorkout(p, "A", "B", 1)

}

func TestTokenize(t *testing.T) {
	expectedLen := func(s string, length int, code parseErrors) {
		tokens, e := tokenize(s)
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

/*
It is possible to parse a matcher from a string.  This is obviously great for user input.  See the parse interface documentation for examples
*/
func ExampleParser_1() {
	//We need to give the parser a context.  In this case it is a single field, of type int
	p, _ := NewParser(int(1))

	match, _ := p.Parse("_ = 1")

	//Confirm that the mater works
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
This is more of a reference of how the parser works for a lone field
*/
func ExampleParser_2() {
	printIt := func(parser Parser, expr string, entity interface{}) {
		match, err := parser.Parse(expr)
		if err != nil {
			fmt.Println("There was an error parsing the expression:", expr)
			return
		}
		result, err := match.Match(entity)
		if err != nil {
			fmt.Printf("There was an error matching entity %v for expression %v\n", entity, expr)
			return
		}
		if result {
			fmt.Printf("Expression '%v' matches '%v'\n", expr, entity)
		} else {
			fmt.Printf("Expression '%v' does not match '%v'\n", expr, entity)
		}

	}
	//We need to give the parser a context.  In this case it is a single field, of type int
	p, _ := NewParser(int(1))

	//Equals
	printIt(p, "_ = 1", int(1))
	printIt(p, "_ = 1", int(2))

	//Not equals
	printIt(p, "_ != 1", int(1))
	printIt(p, "_ != 1", int(2))

	//Other comparison operators
	printIt(p, "_ < 1", int(1))
	printIt(p, "_ <= 1", int(1))
	printIt(p, "_ > 1", int(1))
	printIt(p, "_ >= 1", int(1))

	//Compound expressions
	printIt(p, "_ = 1 AND _ = 2", int(1))
	printIt(p, "_ = 1 OR _ = 2", int(1))

	//The empty parser always matches
	printIt(p, "", int(1))

	//Unparsable expressions
	printIt(p, "_ = \"Bacon\"", int(1)) //Can't parse a string when the context is set to an int
	printIt(p, "_ = \"1\"", int(1))     //String conversion does not happen
	printIt(p, "_ = 1.0", int(1))       //Can't parse a floar when context is set to int

	//Output:
	//Expression '_ = 1' matches '1'
	//Expression '_ = 1' does not match '2'
	//Expression '_ != 1' does not match '1'
	//Expression '_ != 1' matches '2'
	//Expression '_ < 1' does not match '1'
	//Expression '_ <= 1' matches '1'
	//Expression '_ > 1' does not match '1'
	//Expression '_ >= 1' matches '1'
	//Expression '_ = 1 AND _ = 2' does not match '1'
	//Expression '_ = 1 OR _ = 2' matches '1'
	//Expression '' matches '1'
	//There was an error parsing the expression: _ = "Bacon"
	//There was an error parsing the expression: _ = "1"
	//There was an error parsing the expression: _ = 1.0
}

/*
This is more of a reference of how the parser works with a composite field
*/
func ExampleParser_3() {
	printIt := func(parser Parser, expr string, entity interface{}) {
		match, err := parser.Parse(expr)
		if err != nil {
			fmt.Println("There was an error parsing the expression:", expr)
			return
		}
		result, err := match.Match(entity)
		if err != nil {
			fmt.Printf("There was an error matching entity %v for expression %v\n", entity, expr)
			return
		}
		if result {
			fmt.Printf("Expression '%v' matches '%v'\n", expr, entity)
		} else {
			fmt.Printf("Expression '%v' does not match '%v'\n", expr, entity)
		}

	}
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

	printIt(p, "A = 0", Foo{})                                   //A simple equlity test matching the zero object
	printIt(p, "A > 0", Foo{})                                   //A simple greater than not matching
	printIt(p, "A = B", Foo{})                                   //Field equality matching
	printIt(p, "A = B", Foo{A: 1, B: 1})                         //Field equality matching still matching
	printIt(p, "A = B AND Name = \"Bacon\"", Foo{Name: "Bacon"}) //Bacon makes things work :)

	//Associativity is right to left
	printIt(p, "A = 0 AND B = 0 OR Name = \"Bacon\"", Foo{})
	printIt(p, "A = 0 AND B = 0 OR Name = \"Bacon\"", Foo{A: 1, Name: "Bacon"})
	printIt(p, "A = 0 AND B = 0 OR Name = \"Bacon\"", Foo{A: 1})

	//Parens overide associativity
	printIt(p, "A = 0 AND (B = 0 OR Name = \"Bacon\")", Foo{})
	printIt(p, "A = 0 AND (B = 0 OR Name = \"Bacon\")", Foo{B: 1, Name: "Bacon"})
	printIt(p, "A = 0 AND (B = 0 OR Name = \"Bacon\")", Foo{A: 1, Name: "Bacon"})
	printIt(p, "A = 0 AND (B = 0 OR Name = \"Bacon\")", Foo{A: 1, B: 0})

	//Not inverts any predicate
	printIt(p, "A = 0 AND NOT (B = 0 OR Name = \"Bacon\")", Foo{})
	printIt(p, "A = 0 AND NOT (B = 0 OR Name = \"Bacon\")", Foo{B: 1})

	printIt(p, "A = C", Foo{}) //Field equlity not parsing, because A and C are different

	//Output:
	//Expression 'A = 0' matches '{0 0 }'
	//Expression 'A > 0' does not match '{0 0 }'
	//Expression 'A = B' matches '{0 0 }'
	//Expression 'A = B' matches '{1 1 }'
	//Expression 'A = B AND Name = "Bacon"' matches '{0 0 Bacon}'
	//Expression 'A = 0 AND B = 0 OR Name = "Bacon"' matches '{0 0 }'
	//Expression 'A = 0 AND B = 0 OR Name = "Bacon"' matches '{1 0 Bacon}'
	//Expression 'A = 0 AND B = 0 OR Name = "Bacon"' does not match '{1 0 }'
	//Expression 'A = 0 AND (B = 0 OR Name = "Bacon")' matches '{0 0 }'
	//Expression 'A = 0 AND (B = 0 OR Name = "Bacon")' matches '{0 1 Bacon}'
	//Expression 'A = 0 AND (B = 0 OR Name = "Bacon")' does not match '{1 0 Bacon}'
	//Expression 'A = 0 AND (B = 0 OR Name = "Bacon")' does not match '{1 0 }'
	//Expression 'A = 0 AND NOT (B = 0 OR Name = "Bacon")' does not match '{0 0 }'
	//Expression 'A = 0 AND NOT (B = 0 OR Name = "Bacon")' matches '{0 1 }'
	//There was an error parsing the expression: A = C
}

/*
This demonstrates how the parser and pretty printer can be combined to go full circle.  This ability to go full circle breaks down if lambda are used in a yielder or matcher.

Some of the expressions are rewritten to be logically equivalent.  This is an optimization done on the internal data structures for performance reasons, not just part of the pretty printer
*/
func ExampleParser_4() {
	printIt := func(parser Parser, expr string, entity interface{}) {
		match, err := parser.Parse(expr)
		printer := NewDefaultPrinter()
		if err != nil {
			fmt.Println("There was an error parsing the expression:", expr)
			return
		}
		result, err := match.Match(entity)
		printedExpr, _ := printer.Print(match)
		if err != nil {
			fmt.Printf("There was an error matching entity %v for expression %v\n", entity, printedExpr)
			return
		}
		if result {
			fmt.Printf("Expression '%v' matches '%v'\n", printedExpr, entity)
		} else {
			fmt.Printf("Expression '%v' does not match '%v'\n", printedExpr, entity)
		}

	}
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

	//A simple equlity test matching the zero object
	printIt(p, "A = 0", Foo{})

	//Notice that the parents are plucked when they aren't required
	printIt(p, "(A = 0)", Foo{})

	//But the parens are respected when they are
	printIt(p, "A = 0 AND (B = 0 OR Name = \"Bacon\")", Foo{})

	//Expresions that use NOT are rewritten if possible
	printIt(p, "NOT (A = 0)", Foo{})
	printIt(p, "NOT (A != B)", Foo{})

	//But preserved when required
	printIt(p, "NOT (A = 0 AND B =0)", Foo{})

	//Output:
	//Expression 'A = 0' matches '{0 0 }'
	//Expression 'A = 0' matches '{0 0 }'
	//Expression 'A = 0 AND (B = 0 OR Name = "Bacon")' matches '{0 0 }'
	//Expression 'A != 0' does not match '{0 0 }'
	//Expression 'A = B' matches '{0 0 }'
	//Expression 'NOT (A = 0 AND B = 0)' does not match '{0 0 }'
}

/*
This example shows the multiple ways to parse a boolean
*/
func ExampleParser_5() {
	printIt := func(parser Parser, expr string, entity interface{}) {
		match, err := parser.Parse(expr)
		if err != nil {
			fmt.Println("There was an error parsing the expression:", expr)
			return
		}
		result, err := match.Match(entity)
		if err != nil {
			fmt.Printf("There was an error matching entity %v for expression %v\n", entity, expr)
			return
		}
		if result {
			fmt.Printf("Expression '%v' matches '%v'\n", expr, entity)
		} else {
			fmt.Printf("Expression '%v' does not match '%v'\n", expr, entity)
		}

	}
	p, _ := NewParser(true)

	//There are mutliple proper ways to express true an false
	printIt(p, "_ = true", true)
	printIt(p, "_ = false", false)
	printIt(p, "_ = 1", false)
	printIt(p, "_ = 0", false)
	printIt(p, "_ != T", false)
	printIt(p, "_ != F", false)
	//Output:
	//Expression '_ = true' matches 'true'
	//Expression '_ = false' matches 'false'
	//Expression '_ = 1' does not match 'false'
	//Expression '_ = 0' matches 'false'
	//Expression '_ != T' matches 'false'
	//Expression '_ != F' does not match 'false'
}
