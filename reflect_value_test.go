package main

import (
	//"fmt"
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
