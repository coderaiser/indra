package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func parseFile(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestRemoveUnusedImportName(t *testing.T) {
	r := &RemoveUnusedImport{}
	if r.Name() != "remove-unused-import" {
		t.Errorf("unexpected name: %s", r.Name())
	}
}

func TestRemoveUnusedImportUsed(t *testing.T) {
	src := `package p
import "fmt"
func f() { fmt.Println("hi") }
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	results := r.Check(f, nil)
	if len(results) != 0 {
		t.Errorf("expected 0 violations, got %d", len(results))
	}
}

func TestRemoveUnusedImportUnused(t *testing.T) {
	src := `package p
import "fmt"
func f() {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	r := &RemoveUnusedImport{}
	results := r.Check(f, fset)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Message != `remove unused import: "fmt"` {
		t.Errorf("unexpected message: %s", results[0].Message)
	}
}

func TestRemoveUnusedImportFix(t *testing.T) {
	src := `package p
import "fmt"
func f() {}
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	modified := r.Fix(f, nil)
	if !modified {
		t.Error("expected modified=true")
	}
	results := r.Check(f, nil)
	if len(results) != 0 {
		t.Error("expected no violations after fix")
	}
}

func TestRemoveUnusedImportFixLeavesUsed(t *testing.T) {
	src := `package p
import "fmt"
func f() { fmt.Println("hi") }
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	modified := r.Fix(f, nil)
	if modified {
		t.Error("expected modified=false")
	}
}

func TestRemoveUnusedImportBlank(t *testing.T) {
	src := `package p
import _ "fmt"
func f() {}
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	results := r.Check(f, nil)
	if len(results) != 0 {
		t.Errorf("expected 0 violations for blank import, got %d", len(results))
	}
}

func TestRemoveUnusedImportDot(t *testing.T) {
	src := `package p
import . "fmt"
func f() { Println("hi") }
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	results := r.Check(f, nil)
	if len(results) != 0 {
		t.Errorf("expected 0 violations for dot import, got %d", len(results))
	}
}

func TestRemoveUnusedImportAliasUsed(t *testing.T) {
	src := `package p
import foo "fmt"
func f() { foo.Println("hi") }
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	results := r.Check(f, nil)
	if len(results) != 0 {
		t.Errorf("expected 0 violations for used alias, got %d", len(results))
	}
}

func TestRemoveUnusedImportAliasUnused(t *testing.T) {
	src := `package p
import foo "fmt"
func f() {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	r := &RemoveUnusedImport{}
	results := r.Check(f, fset)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation for unused alias, got %d", len(results))
	}
	if results[0].Message != `remove unused import: "fmt"` {
		t.Errorf("unexpected message: %s", results[0].Message)
	}
}

func TestRemoveUnusedImportAllRemoved(t *testing.T) {
	src := `package p
import "fmt"
func f() {}
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	r.Fix(f, nil)
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			t.Error("expected no import declarations after fix")
		}
	}
}

func TestRemoveUnusedImportMultiBlock(t *testing.T) {
	src := `package p
import (
	"fmt"
	"os"
)
func f() { fmt.Println("hi") }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	r := &RemoveUnusedImport{}
	results := r.Check(f, fset)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if !strings.Contains(results[0].Message, "os") {
		t.Errorf("expected message about os, got: %s", results[0].Message)
	}
	modified := r.Fix(f, fset)
	if !modified {
		t.Error("expected modified=true")
	}
	results = r.Check(f, fset)
	if len(results) != 0 {
		t.Errorf("expected 0 violations after fix, got %d", len(results))
	}
}

func TestRemoveUnusedImportMultiBlockAllUnused(t *testing.T) {
	src := `package p
import (
	"fmt"
	"os"
)
func f() {}
`
	f := parseFile(t, src)
	r := &RemoveUnusedImport{}
	r.Fix(f, nil)
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			t.Error("expected no import declarations after fix")
		}
	}
}