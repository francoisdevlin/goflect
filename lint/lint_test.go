package goflect

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

func TestStructTagValidator(t *testing.T) {
	type ForgottenQuote struct {
		Id int `sql:primary`
	}

	results := ValidateType(&ForgottenQuote{}, NewStructList())
	if len(results) != 1 {
		t.Error("Did not get exactly 1 result back")
	}
	if results[0].Error.Code != TAG_ERROR {
		t.Error("Did not get TAG_ERROR back")
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
	fmt.Println(results)
}
