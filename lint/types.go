/*
There is a lot of imformation associated contained in the struct tags that the go compiler can't directly help us with.  This is a collection of tools to parse the struct tags and determine that they are valid.  There is also a tool to pretty print the struct tag information, in order to make is easier to read.

Using goflect-lint

goflect-lint it the tool to use to verify that a program has properly formed annotations.  You can use it like so:

    goflect-lint hello.go

This will output any errors found in hello.go, and exit non-zero.  It will silently exist zero on success.  For infomration on the linter, please read the ValidateType reference

Using goflect-format

goflect-format is the tool to use in order to pretty print the struct tags.  You can use it like so:

    goflect-format hello.go

This will rewrite any struct tags using the pretty formatter.  This isn't quite 100% compatible with go fmt yet, so you may want to run go fmt on the code afterwards.  This is a known item to fix.  For information on the formatter, please read the FormatStructTag reference
*/
package lint

import (
	"go/token"
)

/*
A StructInfo contains the information about struct and field positions
*/
type StructInfo struct {
	token.Position
	FieldPositions map[string]token.Position
}

/*
This is a basic constructor for the StructInfo type
*/
func NewStructInfo() (output StructInfo) {
	output.FieldPositions = make(map[string]token.Position)
	return output
}

/*
A StructList is used by the linter binary to convery important information about the positions of various types in the source code.  The program the linter genrates needs this type.
*/
type StructList struct {
	token.Position
	Structs map[string]StructInfo
}

/*
This is a basic constructor for the StructLit Type
*/
func NewStructList() (output StructList) {
	output.Structs = make(map[string]StructInfo)
	return output
}

/*
A Result contains a set of errors as well as their associate positions in the source code
*/
type Result struct {
	Error    ValidationError
	Position token.Position
}
