package remove_unused_variables

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"
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

func TestCollectDeclaredNamesWithType(t *testing.T) {
	file := parseDirect(t, `package fixture

type MyType struct{}

const x = 1
func f() {}
var y int
`)
	names := collectDeclaredNames(types.Path{Node: file})
	if !names["MyType"] {
		t.Fatal("missing MyType")
	}
	if !names["x"] {
		t.Fatal("missing x")
	}
	if !names["f"] {
		t.Fatal("missing f")
	}
	if !names["y"] {
		t.Fatal("missing y")
	}
}