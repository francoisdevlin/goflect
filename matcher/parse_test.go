package matcher

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseCodes(t *testing.T) {
	render := func(input string, code parseErrors) {
		p := ParseStruct{Fields: map[string]reflect.Kind{
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

	//Promotion Error
	render("A = C", PROMOTION_ERROR)
	render("C = A", PROMOTION_ERROR)
	render("C = \"Fail\"", PROMOTION_ERROR)
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

	//Strings unknown ATM...
	//p = ParseStruct{Fields: map[string]reflect.Kind{
	//fieldName: reflect.String,
	//}}
	//comparisonWorkout(p, "A", "B", 1)

	//p = ParseStruct{Fields: map[string]reflect.Kind{
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
