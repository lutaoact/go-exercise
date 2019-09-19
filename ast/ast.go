package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	parseDir()
}

func parseFile() {
	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
	println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

func parseDir() {
	fset := token.NewFileSet() // positions are relative to fset
	pkgs, err := parser.ParseDir(fset, "/Users/taolu/telisgo/pkg/job", nil, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	//ast.Print(fset, pkgs)
	ast.Print(fset, pkgs["job"].Files["/Users/taolu/telisgo/pkg/job/migratedata.go"].Scope.Objects)
	objects := pkgs["job"].Files["/Users/taolu/telisgo/pkg/job/migratedata.go"].Scope.Objects
	for _, v := range objects {
		if v.Kind == ast.Var {
			fmt.Println("in:", v.Kind, v.Name)
		} else {
			fmt.Println("out:", v.Kind, v.Name)
		}
	}
}
