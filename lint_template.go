package main

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
}
