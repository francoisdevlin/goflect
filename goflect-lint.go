package main

import (
	"bytes"
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/lint"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const TEMPLATE = `package main

import (
	"fmt"
	"git.sevone.com/sdevlin/goflect.git/lint"
	"go/token"
	"os"
)

type tuple struct {
	Record    interface{}
	Positions lint.StructList
}

func getTypes() []tuple {
	output := []tuple{}
	return output
}

func main() {
	types := getTypes()

	errors := make([]lint.Result, 0)
	for _, typ := range types {
		localErrors := lint.ValidateType(typ.Record, typ.Positions)
		for _, err := range localErrors {
			errors = append(errors, err)
		}
	}

	for _, err := range errors {
		fmt.Println(err.Position, err.Error.Code, err.Error.Message)
	}
	if len(errors) > 0 {
		os.Exit(1)
	}
}`

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
			temp := lint.NewStructInfo()
			temp.Position = getPos(cast.Pos())
			for _, field := range cast.Fields.List {
				if len(field.Names) == 1 {
					temp.FieldPositions[field.Names[0].Name] = getPos(field.Pos())
				}
			}
			v.Structs[t.Name.Name] = temp
		}
	}

	return v
}

func posLiteral(pos token.Position) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.Ident{Name: "token.Position"},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Filename"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"" + pos.Filename + "\"",
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Line"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(int64(pos.Line), 10),
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Column"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.FormatInt(int64(pos.Column), 10),
				},
			},
		},
	}
}

func main() {
	filename := os.Args[1]
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		// Whoops!
	}
	importPath := "./" + filepath.Dir(filename)
	targetPackage := file.Name.Name
	structs := StructCandidateVisitor{fset: fset}
	ast.Walk(&structs, file)

	fset = token.NewFileSet()
	lintfile := "__lint_template.go"
	ioutil.WriteFile(lintfile, []byte(TEMPLATE), os.ModePerm)
	file, err = parser.ParseFile(fset, lintfile, nil, 0)
	os.Remove(lintfile)
	if err != nil {
		// Whoops!
	}

	statements := make([]ast.Stmt, 0)

	f := func(name string, targetPackage string) *ast.CompositeLit {
		return &ast.CompositeLit{
			Type: &ast.Ident{Name: "tuple"},
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key: &ast.Ident{Name: "Record"},
					Value: &ast.CompositeLit{
						Type: &ast.Ident{Name: targetPackage + "." + name},
					},
				},
				&ast.KeyValueExpr{
					Key: &ast.Ident{Name: "Positions"},
					Value: &ast.CompositeLit{
						Type: &ast.Ident{Name: "lint.StructList"},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key:   &ast.Ident{Name: "Position"},
								Value: posLiteral(structs.Structs[name].Position),
							},
						},
					},
				},
			},
		}
	}

	for name, _ := range structs.Structs {
		statements = append(statements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{Name: "output"},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{Name: "append"},
					Args: []ast.Expr{
						&ast.Ident{Name: "output"},
						f(name, targetPackage),
					},
				},
			},
		})
	}
	statements = append(statements, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{Name: "output"},
		},
	})

	getTypes := &ast.FuncDecl{
		Name: &ast.Ident{Name: "getTypes"},
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{&ast.Ident{Name: "output"}},
						Type:  &ast.BasicLit{Kind: token.STRING, Value: "[]tuple"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: statements,
		},
	}
	decls := []ast.Decl{}
	added := false
	for _, d := range file.Decls {
		if gen, ok := d.(*ast.GenDecl); ok && gen.Tok == token.IMPORT && !added {
			gen.Specs = append(gen.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(importPath)},
			})
			added = true
		}
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "getTypes" {
			decls = append(decls, getTypes)
		} else {
			decls = append(decls, d)
		}

	}
	file.Decls = decls
	buffer := bytes.NewBufferString("")
	printer.Fprint(buffer, fset, file)

	tempFile := "abcdef__temp.go"
	ioutil.WriteFile(tempFile, buffer.Bytes(), os.ModePerm)
	cmd := exec.Command("/usr/local/go/bin/go", "run", tempFile)
	out, result := cmd.Output()
	fmt.Print(string(out))
	os.Remove(tempFile)
	if result != nil {
		os.Exit(1)
	}
}
