package goflect

import (
	//"fmt"
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
