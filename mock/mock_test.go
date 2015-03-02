package goflect

import (
	"fmt"
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

func TestMockIds(t *testing.T) {
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

func TestMockImmutable(t *testing.T) {
	type Baz struct {
		Id int64 `sql:"primary,autoincrement"`
		A  int64 `sql:"immutable"`
		B  int64
	}
	mocker := MockerStruct{SkipImmutable: true}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{Id: 0, A: 0, B: 1}) {
		t.Error("Immutable Not Working  - immutable", result)
	}
	mocker = MockerStruct{SkipImmutable: false, SkipId: true}
	mocker.Mock(1, &result)
	if (result != Baz{Id: 0, A: 1, B: 1}) {
		t.Error("Immutable Not working - !immutable && skipid")
	}

	mocker = MockerStruct{SkipImmutable: false, SkipId: false}
	mocker.Mock(1, &result)
	if (result != Baz{Id: 1, A: 1, B: 1}) {
		t.Error("Immutable Not working - !immutable && !skipid")
	}
}

func TestMockSubtype(t *testing.T) {
	type EnumerationType int
	type ErrorType uint
	type UnitType float64
	type Message string

	type Baz struct {
		A EnumerationType
		B ErrorType
		C UnitType
		D Message
	}
	mocker := MockerStruct{}
	result := Baz{}
	mocker.Mock(1, &result)
	if (result != Baz{A: 1, B: 1, C: 1.0, D: "1st"}) {
		t.Error("Enumeration Type Not Set", result)
	}

}

func ExampleMockerStruct_Mock_set() {
	type Bar struct {
		AString   string
		AFloat    float64
		AnInteger int64
		ABool     bool
	}
	mocker := MockerStruct{}
	temp := Bar{}
	mocker.Mock(1, &temp)
	fmt.Println(temp.AString, temp.AFloat, temp.AnInteger, temp.ABool)

	mocker.Mock(2, &temp)
	fmt.Println(temp.AString, temp.AFloat, temp.AnInteger, temp.ABool)
	//Output:
	//1st 1 1 true
	//2nd 2 2 true
}

func ExampleMockerStruct_Mock_skipImmutable() {
	type Bar struct {
		AString   string `sql:"immutable"`
		AFloat    float64
		AnInteger int64
		ABool     bool
	}
	mocker := MockerStruct{SkipImmutable: true}
	temp := Bar{AString: "Not Set"}
	mocker.Mock(1, &temp)
	fmt.Println(temp.AString, temp.AFloat, temp.AnInteger, temp.ABool)
	//Output: Not Set 1 1 true
}

func ExampleMockerStruct_Mock_skipId() {
	type Bar struct {
		AString   string
		AFloat    float64
		AnInteger int64 `sql:"autoincrement"`
		ABool     bool
	}
	mocker := MockerStruct{SkipImmutable: true}
	temp := Bar{}
	mocker.Mock(1, &temp)
	fmt.Println(temp.AString, temp.AFloat, temp.AnInteger, temp.ABool)
	//Output: 1st 1 0 true
}
