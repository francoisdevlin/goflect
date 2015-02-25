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
