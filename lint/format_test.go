package lint

import (
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
