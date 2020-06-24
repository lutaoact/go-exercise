package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	expr, err := parser.ParseExpr("shape.color==color.red")
	if err != nil {
		panic(err)
	}

	// Checking if the expression was binary.
	bExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		panic("expr is not a binary expr.")
	}

	// If the operation is not “==”, die.
	if bExpr.Op != token.EQL {
		panic("the op should have been ==.")
	}

	// Left must be a selector expr, meaning follow with a selector which is “dot” in this case.
	left, ok := bExpr.X.(*ast.SelectorExpr)
	if !ok {
		panic("left should have been a selector expr.")
	}

	// Same as above.
	right, ok := bExpr.Y.(*ast.SelectorExpr)
	if !ok {
		panic("right should have been a selector expr.")
	}

	// Checking for attributes.
	if left.Sel.Name != "color" {
		panic("left should have had a color attr.")
	}

	// Same as above.
	if right.Sel.Name != "red" {
		panic("right should have had a red attr.")
	}

	// Then we finally gofmt the code and print it to stdout.
	if err := format.Node(os.Stdout, token.NewFileSet(), expr); err != nil {
		panic(err)
	}

	fmt.Println()

	// 手动构建语法树，生成相同代码
	newExpr := &ast.BinaryExpr{
		X: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "shape",
			},
			Sel: &ast.Ident{
				Name: "color",
			},
		},
		Op: token.EQL,
		Y: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "color",
			},
			Sel: &ast.Ident{
				Name: "red",
			},
		},
	}
	if err := format.Node(os.Stdout, token.NewFileSet(), newExpr); err != nil {
		panic(err)
	}
}
