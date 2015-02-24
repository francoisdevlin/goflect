package goflect

import (
	"go/token"
)

type StructInfo struct {
	token.Position
	FieldPositions map[string]token.Position
}

func NewStructInfo() (output StructInfo) {
	output.FieldPositions = make(map[string]token.Position)
	return output
}

type StructList struct {
	token.Position
	fset    *token.FileSet
	Structs map[string]StructInfo
}

func NewStructList() (output StructList) {
	output.Structs = make(map[string]StructInfo)
	return output
}
