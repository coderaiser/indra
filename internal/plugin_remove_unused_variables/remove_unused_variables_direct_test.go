package remove_unused_variables

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseDirect(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

func TestImportFindingEnd(t *testing.T) {
	file := parseDirect(t, `package fixture

import "fmt"

func f() {}
`)
	imports := collectImports(file)
	finding := &importFinding{file: file, spec: imports[0].spec}
	if finding.End() != imports[0].spec.End() {
		t.Fatal("End mismatch")
	}
}

func TestFuncDeclFindingEnd(t *testing.T) {
	file := parseDirect(t, `package fixture

func unusedHelper() {}
`)
	funcDecl := file.Decls[0].(*ast.FuncDecl)
	finding := &funcDeclFinding{file: file, decl: funcDecl}
	if finding.End() != funcDecl.End() {
		t.Fatal("End mismatch")
	}
	if finding.Pos() != funcDecl.Pos() {
		t.Fatal("Pos mismatch")
	}
}
