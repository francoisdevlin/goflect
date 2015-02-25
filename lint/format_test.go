package lint

import (
	"fmt"
	"go/token"
	"testing"
)

func TestFormatString(t *testing.T) {
	assertMatch := func(expected, message string, position token.Position) {
		output, err := FormatStructTag(position, message)
		if err != nil && len(err) > 0 {
			t.Errorf("There was an unexpected error with the input: %v", message)
		}
		if output != expected {
			t.Errorf("There was a formatting error, expected '%v' got '%v'", expected, output)
		}
	}
	assertMatch("sql:\"primary\"", "`sql:\"primary\"`", token.Position{})
	//Cover no backquotes
	assertMatch("sql:\"primary\"", "sql:\"primary\"", token.Position{})
	//Prune missing field
	//assertMatch("", "`sql:\"\"`", token.Position{})

	//Columns are 1 based, hence the off by 1 nature of the tests
	//The first line never gets space appended
	assertMatch("sql:\"primary\"", "`sql:\"primary\"`", token.Position{Column: 1})
	assertMatch("sql:\"primary\"", "`sql:\"primary\"`", token.Position{Column: 2})
	//Description comes first, padded tab
	assertMatch("desc:\"Bacon\"\n\t sql:\"primary\"", "`sql:\"primary\" desc:\"Bacon\"`", token.Position{Column: 2})
}

/*
The sql fields will always be placed in a specific order, which is determined in the goflect package
*/
func ExampleFormatStructTag_sqlReorder() {
	garbledTag := `sql:"not-null,autoincrement,unique,primary,index,nominal,immutable"`

	output, _ := FormatStructTag(token.Position{}, garbledTag)
	fmt.Println(output)
	//Output:
	//sql:"primary, autoincrement, unique, immutable, nominal, not-null, index"
}

/*
The ui fields will always be placed in a specific order, which is determined in the goflect package
*/
func ExampleFormatStructTag_uiReorder() {
	garbledTag := `ui:"redacted,hidden"`

	output, _ := FormatStructTag(token.Position{}, garbledTag)
	fmt.Println(output)
	//Output:
	//ui:"hidden, redacted"
}

/*
The struct tags will be presented in a specific order, so that there is a standard way of documenting these items.  This will help us scale as the number of annotated structs approached the hundreds.

The leading tab is an artifact to have the tags look pretty when combined with go fmt.
*/
func ExampleFormatStructTag_tagReorder() {
	garbledTag := `ui:"hidden" sql"primary" desc:"This is a primary key" valid:"_ >= 0"`

	output, _ := FormatStructTag(token.Position{}, garbledTag)
	fmt.Println(output)
	//Output:
	//desc:"This is a primary key"
	//	valid:"_ >= 0"
	//	sql:"primary"
	//	ui:"hidden"
}
