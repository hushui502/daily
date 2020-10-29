package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	srcPath := os.Args[1]
	fmt.Printf("Parsing source file %s...\n", srcPath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcPath, nil, 0)

	if err != nil {
		panic(err)
	}

	fmt.Println("Found imports:")
	for _, s := range f.Imports {
		fmt.Println(s.Path.Value)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			ast.Print(fset, x.Fun)
		}
		return true
	})

	hasPrint := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			selexpr, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := selexpr.X.(*ast.Ident)
			if !ok || ident.Name != "fmt" {
				return true
			}
			if selexpr.Sel.Name == "Printf" || selexpr.Sel.Name == "Println" {
				// convert compact token position to raw source position for display
				pos := fset.Position(selexpr.Sel.Pos())
				fmt.Printf("Use of `fmt.%s` detected at %v\n", selexpr.Sel.Name, pos)
				hasPrint = true
			}
		}
		return true
	})
	if hasPrint {
		os.Exit(1)
	} else {
		fmt.Println("All good!")
	}
}
