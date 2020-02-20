package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

// identsToFmtStmt generates fmt.Printf statements from identifiers
func identsToFmtStmt(identToExpr map[string]ast.Expr) []ast.Stmt {
	fmtStatements := make([]ast.Stmt, 0)
	for name, expr := range identToExpr {
		fmtStatements = append(
			fmtStatements,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("fmt"),
						Sel: ast.NewIdent("Printf"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"%s = %+v\\n\"",
						},
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf("\"%s\"", name),
						},
						expr,
					},
				},
			},
		)
	}
	return fmtStatements
}

// parseAssigmentStmt parses := or = assigments
//
// Grammar Rules:
// ShortVarDecl = IdentifierList ":=" ExpressionList .
func parseAssigmentStmt(assignStmt *ast.AssignStmt) []ast.Stmt {
	leftIdents := make(map[string]ast.Expr, len(assignStmt.Lhs))
	for _, identExpr := range assignStmt.Lhs {
		switch v := identExpr.(type) {
		case *ast.Ident:
			leftIdents[v.Name] = v
		case *ast.SelectorExpr:
			leftIdents[v.Sel.Name] = v
		}

	}

	return identsToFmtStmt(leftIdents)
}

// parseDeclSpecs parses var or const block declarations
//
// Grammar Rules:
// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
// VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
func parseDeclSpecs(specs []ast.Spec) []ast.Stmt {
	leftIdents := make(map[string]ast.Expr)
	for _, spec := range specs {
		valueSpec := spec.(*ast.ValueSpec)

		for _, ident := range valueSpec.Names {
			leftIdents[ident.Name] = ident
		}
	}

	return identsToFmtStmt(leftIdents)
}

func insertAssigmentDebug(file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		blockStmt, ok := node.(*ast.BlockStmt)
		if !ok {
			return true
		}

		updatedBlockList := make([]ast.Stmt, 0)
		for _, stmt := range blockStmt.List {
			updatedBlockList = append(updatedBlockList, stmt)
			switch v := stmt.(type) {
			case *ast.AssignStmt:
				updatedBlockList = append(updatedBlockList, parseAssigmentStmt(v)...)
			case *ast.DeclStmt:
				decl := v.Decl.(*ast.GenDecl)
				if decl.Tok == token.VAR || decl.Tok == token.CONST {
					updatedBlockList = append(updatedBlockList, parseDeclSpecs(decl.Specs)...)
				}
			}
		}

		blockStmt.List = updatedBlockList
		return true
	})
}

func main() {
	if len(os.Args) != 2 {
		return
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		fmt.Printf("Errors in %s\n", os.Args[1])
		return
	}

	insertAssigmentDebug(file)
	if format.Node(os.Stdout, fset, file) != nil {
		fmt.Printf("Formatter error: %v\n", err)
	}
}
