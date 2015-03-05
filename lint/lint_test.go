package lint

import (
	"fmt"
	"testing"
)

func TestNominalValidator(t *testing.T) {
	type NominalIntMisconfigure struct {
		Id   int `sql:"primary"`
		Name int `sql:"nominal,unique"`
	}

	results := ValidateType(&NominalIntMisconfigure{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != NOMINAL_MISMATCH {
		t.Error("Did not get NOMINAL_MISMATCH back")
	}

	type NominalUniqueMisconfigure struct {
		Id   int    `sql:"primary"`
		Name string `sql:"nominal"`
	}

	results = ValidateType(&NominalUniqueMisconfigure{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != NOMINAL_MISMATCH {
		t.Error("Did not get NOMINAL_MISMATCH back")
	}
}

/*
This demonstrates some of the requirements around the nominal keyword.  The nominal field must be type string.  It must be with something that is able to unique.  The type must have a primary key, and there can only be one nominal field for a given type.  A string primary key may also be nominal.
*/
func ExampleValidateType_nominalConstraints() {
	type NominalIntMisconfigure struct {
		Id   int `sql:"primary"`
		Name int `sql:"nominal,unique"`
	}

	results := ValidateType(&NominalIntMisconfigure{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type NominalUniqueMisconfigure struct {
		Id   int    `sql:"primary"`
		Name string `sql:"nominal"`
	}

	results = ValidateType(&NominalUniqueMisconfigure{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type NominalRepeatedMisconfigure struct {
		Id    int    `sql:"primary"`
		Name  string `sql:"nominal,unique"`
		Value string `sql:"nominal,unique"`
	}

	results = ValidateType(&NominalRepeatedMisconfigure{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type NominalMissingPrimary struct {
		Id   int
		Name string `sql:"nominal,unique"`
	}

	results = ValidateType(&NominalMissingPrimary{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)
	//Output:
	//NOMINAL_MISMATCH Field Name is marked nominal, but is kind int with field "Name"
	//NOMINAL_MISMATCH Field Name is marked nominal, but is not unique with field "Name"
	//NOMINAL_MISCOUNT There can be only one nominal field, but the following are marked, [Name Value] on type "NominalRepeatedMisconfigure"
	//NOMINAL_MISMATCH There is a nominal field without a primary key on type "NominalMissingPrimary"

}

func TestStructTagValidator(t *testing.T) {
	type ForgottenQuote struct {
		Id int `sql:primary`
	}

	results := ValidateType(&ForgottenQuote{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Errorf("Did not get TAG_ERROR back")
	}

	type ExtraTag struct {
		Id int `sql:"primary" bacon`
	}
	results = ValidateType(&ExtraTag{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Error("Did not get TAG_ERROR back")
	}

	type RepeatedTag struct {
		Id int `sql:"primary" sql:"primary"`
	}
	results = ValidateType(&RepeatedTag{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Error("Did not get TAG_ERROR back")
	}

	type ExtraSpace struct {
		Id int `sql: "primary"`
	}
	results = ValidateType(&ExtraSpace{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Error("Did not get TAG_ERROR back")
	}

	type NonExistantFlag struct {
		Id int `sql:"primary,bacon"`
	}
	results = ValidateType(&NonExistantFlag{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Error("Did not get TAG_ERROR back")
	}
}

/*
These are a collection of examples showing various struct tags that don't parse properly
*/
func ExampleValidateType_parseErrors() {
	type ForgottenQuote struct {
		Id int `sql:primary`
	}

	results := ValidateType(&ForgottenQuote{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type ExtraTag struct {
		Id int `sql:"primary" bacon`
	}
	results = ValidateType(&ExtraTag{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type RepeatedTag struct {
		Id int `sql:"primary" sql:"primary"`
	}
	results = ValidateType(&RepeatedTag{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)

	type NonExistantFlag struct {
		Id int `sql:"primary,bacon"`
	}
	results = ValidateType(&NonExistantFlag{}, NewStructList())
	fmt.Println(results[0].Error.Code, results[0].Error.Message)
	//Output:
	//TAG_PARSE_ERROR :primary with field "Id"
	//TAG_PARSE_ERROR There are the wrong number of tokens present with field "Id"
	//TAG_PARSE_ERROR Key sql has been repeated with field "Id"
	//TAG_PARSE_ERROR Flag 'bacon' is not allowed for tag "sql" with field "Id"
}

func TestPrimaryOnceValidator(t *testing.T) {
	type DoublePrimary struct {
		A string `sql:"primary"`
		B string `sql:"primary"`
	}

	results := ValidateType(&DoublePrimary{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != PRIMARY_MISCOUNT {
		t.Error("Did not get PRIMARY_MISCOUNT back")
	}
}

/*
This demonstrates the constraints around a primary key. It must be a type that can be unique (not float, bool or byte[]), and there can only be one for a given type
*/
func ExampleValidateType_primaryConstraints() {
	printErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				fmt.Println(err.Error.Code, err.Error.Message)
			}
		}
	}

	type RepeatedPrimary struct {
		Id    int `sql:"primary"`
		Value int `sql:"primary"`
	}
	type BooleanPrimary struct {
		Id bool `sql:"primary"`
	}
	type FloatPrimary struct {
		Id float64 `sql:"primary"`
	}
	printErrors(
		&RepeatedPrimary{},
		&BooleanPrimary{},
		&FloatPrimary{},
	)
	//Output:
	//PRIMARY_MISCOUNT There can be only one primary field, but the following are marked, [Id Value] on type "RepeatedPrimary"
	//PRIMARY_MISMATCH Field Id is marked unique, but is kind bool with field "Id"
	//PRIMARY_MISMATCH Field Id is marked unique, but is kind float64 with field "Id"
}

func TestNominalOnceValidator(t *testing.T) {
	type DoubleNominal struct {
		Id int    `sql:"primary"`
		A  string `sql:"nominal,unique"`
		B  string `sql:"nominal,unique"`
	}

	results := ValidateType(&DoubleNominal{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != NOMINAL_MISCOUNT {
		t.Error("Did not get NOMINAL_MISCOUNT back")
	}
}

/*
This demonstrates the constraints around an autoinc annoation.  It must be a primary key alreaady, and something int-like
*/
func ExampleValidateType_autoincConstraints() {
	printErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				fmt.Println(err.Error.Code, err.Error.Message)
			}
		}
	}

	type NoPrimary struct {
		Id string `sql:"autoincrement"`
	}
	type StringAutoinc struct {
		Id string `sql:"primary,autoincrement"`
	}
	printErrors(
		&NoPrimary{},
		&StringAutoinc{},
	)
	//Output:
	//AUTOINC_ERROR Marked autoinc, but not primary with field "Id"
	//AUTOINC_ERROR Field is marked autoinc, but is kind string with field "Id"
}

/*
This demonstrates errors for defaults that cannon be parsed to the appropriate kind
*/
func ExampleValidateType_defaultMistmatches() {
	printErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				fmt.Println(err.Error.Code, err.Error.Message)
			}
		}
	}

	type BoolStringMismatch struct {
		Value bool `default:"fail"`
	}
	type FloatStringMismatch struct {
		Value float64 `default:"fail"`
	}
	type IntStringMismatch struct {
		Value int `default:"fail"`
	}
	type IntFloatMismatch struct {
		Value int `default:"10.1"`
	}
	type UintStringMismatch struct {
		Value uint `default:"fail"`
	}
	type UintFloatMismatch struct {
		Value uint `default:"10.1"`
	}
	type UintIntMismatch struct {
		Value uint `default:"-1"`
	}
	printErrors(
		&BoolStringMismatch{},
		&FloatStringMismatch{},
		&IntStringMismatch{},
		&IntFloatMismatch{},
		&UintStringMismatch{},
		&UintFloatMismatch{},
		&UintIntMismatch{},
	)
	//Output:
	//BAD_DEFAULT_VALUE Unable to convert "fail" to kind bool with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "fail" to kind float64 with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "fail" to kind int with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "10.1" to kind int with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "fail" to kind uint with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "10.1" to kind uint with field "Value"
	//BAD_DEFAULT_VALUE Unable to convert "-1" to kind uint with field "Value"
}

/*
This demonstrated some of the errors that can occur if the validator expression is not parsable
*/
func ExampleValidateType_validExprParseErrors() {
	printErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				fmt.Println(err.Error.Code, err.Error.Message)
			}
		}
	}

	type MismatchEquality struct {
		A int `valid:"A=B"`
		B int64
	}
	type IncompleteExpression struct {
		A int `valid:"A= "`
		B int
	}
	type MismatchedParen struct {
		A int `valid:"(A=B "`
		B int
	}
	type DanglingParen struct {
		A int `valid:"(A=B))"`
		B int
	}
	printErrors(
		&MismatchEquality{},
		&IncompleteExpression{},
		&MismatchedParen{},
		&DanglingParen{},
		//Can't compare ints an strings
		struct {
			A int `valid:"A=\"Bacon\""`
		}{},
		//Operators must bu a known set
		struct {
			A int `valid:"A==1"`
		}{},
	)
	//Output:
	//VALIDATOR_PARSE_ERROR Cannot compare fields A and B, they are different kinds on type "MismatchEquality"
	//VALIDATOR_PARSE_ERROR The message has a trailing entry: = on type "IncompleteExpression"
	//VALIDATOR_PARSE_ERROR There is an leading paren without its mate on type "MismatchedParen"
	//VALIDATOR_PARSE_ERROR Unknown Field provided: ) on type "DanglingParen"
	//VALIDATOR_PARSE_ERROR Could not promote field A to kind int for value '"Bacon"' on type ""
	//VALIDATOR_PARSE_ERROR Operation type is not supported: == on type ""
}

func TestErrorCodeSerialization(t *testing.T) {
	codes := map[ErrorCode]string{
		NOMINAL_MISMATCH:      "NOMINAL_MISMATCH",
		PRIMARY_MISMATCH:      "PRIMARY_MISMATCH",
		TAG_ERROR:             "TAG_PARSE_ERROR",
		NOMINAL_MISCOUNT:      "NOMINAL_MISCOUNT",
		PRIMARY_MISCOUNT:      "PRIMARY_MISCOUNT",
		AUTOINC_ERROR:         "AUTOINC_ERROR",
		UNIQUE_ERROR:          "UNIQUE_ERROR",
		BAD_DEFAULT:           "BAD_DEFAULT_VALUE",
		VALIDATOR_PARSE_ERROR: "VALIDATOR_PARSE_ERROR",
		//Testing the does not exist case for completeness
		ErrorCode(-1): "",
	}

	for code, str := range codes {
		if code.String() != str {
			t.Errorf("Code %v string error.  Expected: %v, got %v", code, str, code.String())
		}
	}
}

func TestValidAutoincStructure(t *testing.T) {
	noErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				t.Errorf("Found error, code: '%v', message %v", err.Error.Code, err.Error.Message)
			}
		}
	}

	noErrors(
		struct {
			A int `sql:"autoincrement,primary"`
		}{},
		struct {
			A int64 `sql:"autoincrement,primary"`
		}{},
		struct {
			A int32 `sql:"autoincrement,primary"`
		}{},
		struct {
			A int16 `sql:"autoincrement,primary"`
		}{},
		struct {
			A int8 `sql:"autoincrement,primary"`
		}{},
		struct {
			A uint `sql:"autoincrement,primary"`
		}{},
		struct {
			A uint64 `sql:"autoincrement,primary"`
		}{},
		struct {
			A uint32 `sql:"autoincrement,primary"`
		}{},
		struct {
			A uint16 `sql:"autoincrement,primary"`
		}{},
		struct {
			A uint8 `sql:"autoincrement,primary"`
		}{},
	)
}

func TestValidDefaultStructure(t *testing.T) {
	noErrors := func(args ...interface{}) {
		for _, arg := range args {
			results := ValidateType(arg, NewStructList())
			for _, err := range results {
				t.Errorf("Found error, code: '%v', message %v", err.Error.Code, err.Error.Message)
			}
		}
	}

	noErrors(
		struct {
			A int `default:"1"`
		}{},
		struct {
			A int64 `default:"1"`
		}{},
		struct {
			A int32 `default:"1"`
		}{},
		struct {
			A int16 `default:"1"`
		}{},
		struct {
			A int8 `default:"1"`
		}{},
		struct {
			A uint `default:"1"`
		}{},
		struct {
			A uint64 `default:"1"`
		}{},
		struct {
			A uint32 `default:"1"`
		}{},
		struct {
			A uint16 `default:"1"`
		}{},
		struct {
			A uint8 `default:"1"`
		}{},
		struct {
			A float64 `default:"1"`
		}{},
		struct {
			A float32 `default:"1"`
		}{},
		struct {
			A string `default:"1"`
		}{},
	)
}
