package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	fd, err := os.Open(dir)
	fmt.Println(err)
	defer fd.Close()

	list, err := fd.Readdir(-1)
	fmt.Println(err)
	for _, v := range list {
		fmt.Println(v.Name())
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, pkgs["job"].Files["/Users/taolu/telisgo/pkg/job/migratedata.go"].Scope.Objects)
}
