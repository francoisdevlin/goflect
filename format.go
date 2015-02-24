package main

import (
	"bytes"
	//"bufio"
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/lint"
	"go/ast"
	//"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
)

type StructCandidateVisitor struct {
	fset *token.FileSet
	lint.StructList
}

func (v *StructCandidateVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if v.Structs == nil {
		v.Structs = make(map[string]lint.StructInfo)
	}
	getPos := func(p token.Pos) token.Position {
		position := v.fset.File(p).Position(p)
		return position
	}
	switch t := node.(type) {
	case *ast.TypeSpec:
		cast, ok := t.Type.(*ast.StructType)
		if ok {
			for _, field := range cast.Fields.List {
				if field.Tag != nil {
					pos := getPos(field.Tag.ValuePos)
					tagString := field.Tag.Value
					newString, errors := lint.FormatStructTag(pos, tagString)
					if len(errors) == 0 {
						newString = "`" + newString + "`"
						field.Tag = &ast.BasicLit{Kind: token.STRING, Value: newString, ValuePos: field.Tag.ValuePos}
					}
				}
			}
		}
	}

	return v
}

func main() {
	rest := os.Args[2:]
	for _, filename := range rest {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("No such file or directory: %s\n", filename)
			os.Exit(1)
		}
	}

	for _, filename := range rest {
		buffer := bytes.NewBufferString("")

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			// Whoops!
		}
		structs := StructCandidateVisitor{fset: fset}
		ast.Walk(&structs, file)
		printer.Fprint(buffer, fset, file)
		ioutil.WriteFile(filename, buffer.Bytes(), os.ModePerm)
	}
}
