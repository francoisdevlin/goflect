package goflect

import (
	"fmt"
	"testing"
)

func TestSqlFields(t *testing.T) {
	compare := func(expected SqlInfo, record interface{}) {
		records := GetInfo(record)
		result := records[0].SqlInfo
		if result != expected {
			t.Errorf("Error, expects %v, got %v", expected, result)
		}
	}

	type T01 struct {
		Id int `sql:"primary,autoincrement"`
	}
	compare(SqlInfo{
		IsPrimary:       true,
		IsAutoincrement: true,
		IsImmutable:     true,
		IsUnique:        true,
		IsIndexed:       true,
	}, &T01{})

	type T02 struct {
		Id int `sql:"primary"`
	}
	compare(SqlInfo{
		IsPrimary: true,
		IsUnique:  true,
		IsIndexed: true,
	}, &T02{})

	type T03 struct {
		Id int `sql:"immutable"`
	}
	compare(SqlInfo{
		IsNullable:  true,
		IsImmutable: true,
	}, &T03{})

	type T04 struct {
		Id int `sql:"index"`
	}
	compare(SqlInfo{
		IsNullable: true,
		IsIndexed:  true,
	}, &T04{})

	type T05 struct {
		Id int `sql:"unique"`
	}
	compare(SqlInfo{
		IsUnique:  true,
		IsIndexed: true,
	}, &T05{})

	type T06 struct {
		Id int `sql:"not-null"`
	}
	compare(SqlInfo{}, &T06{})

	type T07 struct {
		Id int `sql:"nominal"`
	}
	compare(SqlInfo{
		IsNullable: true,
		IsNominal:  true,
	}, &T07{})
}

func TestValidatorFields(t *testing.T) {
	compare := func(expected ValidatorInfo, record interface{}) {
		records := GetInfo(record)
		result := records[0].ValidatorInfo
		if result.IsRequired != expected.IsRequired ||
			result.MaxValue != expected.MaxValue ||
			result.MinValue != expected.MinValue ||
			result.MatchRegex != expected.MatchRegex ||
			len(result.InEnum) != len(expected.InEnum) ||
			false {
			t.Errorf("Error, expects %v, got %v", expected, result)
		}
	}

	type T01 struct {
		Id int `valid:"required" valid-enum:"A,B,C" valid-max:"10.0" valid-min:"1.0" valid-regex:"\\d+"`
	}
	compare(ValidatorInfo{
		IsRequired: true,
		MaxValue:   "10.0",
		MinValue:   "1.0",
		MatchRegex: "\\d+",
		InEnum:     []string{"A", "B", "C"},
	}, &T01{})
}

func ExampleReflectValue_GetFieldInfo() {
	type Bar struct {
		Id int64
	}

	info := GetInfo(&Bar{})
	fmt.Println(info[0].Name, info[0].Kind)
	//Output: Id int64
}

func ExampleReflectValue_GetFieldSqlInfo_primary() {
	type Bar struct {
		Id int64 `sql:"primary,autoincrement"`
	}

	info := GetInfo(&Bar{})
	fmt.Println(info[0].IsPrimary, info[0].IsAutoincrement)
	//Output: true true
}

func ExampleReflectValue_GetFieldSqlInfo_immutable() {
	type Bar struct {
		Id int64 `sql:"immutable"`
	}

	info := GetInfo(&Bar{})
	fmt.Println(info[0].IsImmutable)
	//Output: true
}

func ExampleReflectValue_GetFieldUiInfo() {
	type Bar struct {
		Id int64 `desc:"An Id Field" default:"1" ui:"hidden"`
	}

	info := GetInfo(&Bar{})
	fmt.Println(info[0].Description, info[0].Default, info[0].Hidden)
	//Output: An Id Field 1 true
}
