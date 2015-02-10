package main

import (
	//"database/sql"
	//"fmt"
	//_ "github.com/mattn/go-sqlite3"
	//"reflect"
	"testing"
)

func TestMockInts(t *testing.T) {
	type Baz struct {
		I   int
		I64 int64
		I32 int32
		I16 int16
		I8  int8
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{I: 1, I64: 1, I32: 1, I16: 1, I8: 1}) {
		t.Error("Signed Integers not working")
	}
}

func TestMockUints(t *testing.T) {
	type Baz struct {
		U   int
		U64 int64
		U32 int32
		U16 int16
		U8  int8
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{U: 1, U64: 1, U32: 1, U16: 1, U8: 1}) {
		t.Error("Unsigned Integers not working")
	}
}

func TestMockBools(t *testing.T) {
	type Baz struct {
		B bool
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{B: true}) {
		t.Error("Bools not working")
	}
}

func TestMockStrings(t *testing.T) {
	type Baz struct {
		S string
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{S: "1st"}) {
		t.Error("String not working")
	}
	mocker.Mock(2, &result)
	if (result != Baz{S: "2nd"}) {
		t.Error("String not working")
	}
	mocker.Mock(3, &result)
	if (result != Baz{S: "3rd"}) {
		t.Error("String not working")
	}
	mocker.Mock(4, &result)
	if (result != Baz{S: "4th"}) {
		t.Error("String not working")
	}
}

func TestMockFloats(t *testing.T) {
	type Baz struct {
		F32 float32
		F64 float64
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{F32: 1.0, F64: 1.0}) {
		t.Error("Floats not working")
	}
}

func TestMockctsIds(t *testing.T) {
	type Baz struct {
		Id int64 `sql:"primary,autoincrement"`
	}
	mocker := MockerStruct{SkipId: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{}) {
		t.Error("SkipId not working - true case")
	}
	mocker = MockerStruct{SkipId: false}
	mocker.Mock(1, &result)
	if (result != Baz{Id: 1}) {
		t.Error("SkipId not working - false case")
	}
}