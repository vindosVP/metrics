package mainexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mainexit",
	Doc:  "check of os.exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			isOsExit := func(n ast.Node) {
				if x, ok := n.(*ast.CallExpr); ok {
					if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
						pkg, isIdent := fun.X.(*ast.Ident)
						if isIdent && pkg.Name == "os" && fun.Sel.Name == "Exit" {
							pass.Reportf(x.Pos(), "os.Exit call in main function")
						}
					}
				}
			}

			isMainFunc := func(n ast.Node) {
				if x, ok := n.(*ast.FuncDecl); ok {
					if x.Name.Name == "main" {
						for _, v := range x.Body.List {
							if stmt, ok := v.(*ast.ExprStmt); ok {
								isOsExit(stmt.X)
							}
						}
					}
				}
			}

			switch x := n.(type) {
			case *ast.File:
				if x.Name.Name == "main" {
					for _, v := range x.Decls {
						isMainFunc(v)
					}
				}
			}

			return true
		})
	}
	return nil, nil
}
