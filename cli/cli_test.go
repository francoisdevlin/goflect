package goflect

import (
	"fmt"
	"testing"
)

func TestParseIntArgs(t *testing.T) {
	type NoDefaults struct {
		I   int
		I64 int64
		I32 int32
		I16 int16
		I8  int8
	}
	type Defaults struct {
		I   int   `default:"1"`
		I64 int64 `default:"64"`
		I32 int32 `default:"32"`
		I16 int16 `default:"16"`
		I8  int8  `default:"8"`
	}

	noDefs := NoDefaults{}
	FlagSetup(&noDefs, []string{"Name"})
	if (noDefs != NoDefaults{}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-I", "100", "-I64", "164", "-I32", "132", "-I16", "116", "-I8", "108"})
	if (noDefs != NoDefaults{100, 164, 132, 116, 108}) {
		t.Error("No Default Parsing Error")
	}

	defs := Defaults{}
	FlagSetup(&defs, []string{"Name"})
	if (defs != Defaults{1, 64, 32, 16, 8}) {
		t.Error("Default Parsing Error")
	}

	defs = Defaults{}
	FlagSetup(&defs, []string{"Name", "-I", "100", "-I64", "164", "-I32", "132", "-I16", "116", "-I8", "108"})
	if (defs != Defaults{100, 164, 132, 116, 108}) {
		t.Error("Default Parsing Error")
	}
}

func TestParseUintArgs(t *testing.T) {
	type NoDefaults struct {
		U   uint
		U64 uint64
		U32 uint32
		U16 uint16
		U8  uint8
	}
	type Defaults struct {
		U   uint   `default:"1"`
		U64 uint64 `default:"64"`
		U32 uint32 `default:"32"`
		U16 uint16 `default:"16"`
		U8  uint8  `default:"8"`
	}

	noDefs := NoDefaults{}
	FlagSetup(&noDefs, []string{"Name"})
	if (noDefs != NoDefaults{}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-U", "100", "-U64", "164", "-U32", "132", "-U16", "116", "-U8", "108"})
	if (noDefs != NoDefaults{100, 164, 132, 116, 108}) {
		t.Error("No Default Parsing Error")
	}

	defs := Defaults{}
	FlagSetup(&defs, []string{"Name"})
	if (defs != Defaults{1, 64, 32, 16, 8}) {
		t.Error("Default Parsing Error")
	}

	defs = Defaults{}
	FlagSetup(&defs, []string{"Name", "-U", "100", "-U64", "164", "-U32", "132", "-U16", "116", "-U8", "108"})
	if (defs != Defaults{100, 164, 132, 116, 108}) {
		t.Error("Default Parsing Error")
	}
}

func TestParseFloatArgs(t *testing.T) {
	type NoDefaults struct {
		F64 float64
		F32 float32
	}
	type Defaults struct {
		F64 float64 `default:"64"`
		F32 float32 `default:"32"`
	}

	noDefs := NoDefaults{}
	FlagSetup(&noDefs, []string{"Name"})
	if (noDefs != NoDefaults{}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-F64", "164", "-F32", "132"})
	if (noDefs != NoDefaults{164, 132}) {
		t.Error("No Default Parsing Error")
	}

	defs := Defaults{}
	FlagSetup(&defs, []string{"Name"})
	if (defs != Defaults{64, 32}) {
		t.Error("Default Parsing Error")
	}

	defs = Defaults{}
	FlagSetup(&defs, []string{"Name", "-F64", "164", "-F32", "132"})
	if (defs != Defaults{164, 132}) {
		t.Error("Default Parsing Error")
	}
}

func TestParseBoolArgs(t *testing.T) {
	type NoDefaults struct {
		B bool
	}

	type DefaultTrue struct {
		B bool `default:"true"`
	}

	noDefs := NoDefaults{}
	FlagSetup(&noDefs, []string{"Name"})
	if (noDefs != NoDefaults{}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-B", "true"})
	if (noDefs != NoDefaults{true}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-B"})
	if (noDefs != NoDefaults{true}) {
		t.Error("No Default Parsing Error")
	}

	defTrue := DefaultTrue{}
	FlagSetup(&defTrue, []string{"Name"})
	if (defTrue != DefaultTrue{true}) {
		t.Error("Default Parsing Error")
	}

	defTrue = DefaultTrue{}
	FlagSetup(&defTrue, []string{"Name", "-B=false"})
	if (defTrue != DefaultTrue{}) {
		t.Error("Default Parsing Error")
	}
}

func TestParseStringArgs(t *testing.T) {
	type NoDefaults struct {
		S string
	}

	type Defaults struct {
		S string `default:"Awesome"`
	}

	noDefs := NoDefaults{}
	FlagSetup(&noDefs, []string{"Name"})
	if (noDefs != NoDefaults{}) {
		t.Error("No Default Parsing Error")
	}

	noDefs = NoDefaults{}
	FlagSetup(&noDefs, []string{"Name", "-S", "Awesome"})
	if (noDefs != NoDefaults{"Awesome"}) {
		t.Error("No Default Parsing Error")
	}

	defTrue := Defaults{}
	FlagSetup(&defTrue, []string{"Name"})
	if (defTrue != Defaults{"Awesome"}) {
		t.Errorf("Default Parsing Error, %v", defTrue)
	}

	defTrue = Defaults{}
	FlagSetup(&defTrue, []string{"Name", "-S", "Not Awesome"})
	if (defTrue != Defaults{"Not Awesome"}) {
		t.Error("Default Parsing Error")
	}
}

func ExampleFlagSetup_basic() {
	type Bar struct {
		//The desc metadata is used for command line help
		Name  string `desc:"This is the name"`
		Value int64  `desc:"This is a value"`
	}
	temp := Bar{}
	FlagSetup(&temp, []string{"AppName", "-Name", "Sean", "-Value", "1"})
	fmt.Println(temp.Name, temp.Value)
	//Output: Sean 1
}

func ExampleFlagSetup_withDefaults() {
	type Bar struct {
		Name  string `default:"Bacon"`
		Value int64  `default:"10"`
	}
	temp := Bar{}
	FlagSetup(&temp, []string{"AppName"})
	fmt.Println(temp.Name, temp.Value)
	//Output: Bacon 10
}
