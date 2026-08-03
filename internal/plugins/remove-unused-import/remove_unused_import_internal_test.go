package remove_unused_import

import (
	"go/ast"
	"go/token"
	"testing"

	"coderaiser/indra/types"
)

func TestPluginReportMessage(t *testing.T) {
	if got := Report(); got != "remove unused import" {
		t.Fatalf("unexpected report message: %q", got)
	}
}

// TestSelfReportMessage verifies the self wrapper forwards the report message,
// covering the method the engine-loader reflects on.
func TestSelfReportMessage(t *testing.T) {
	if got := Self.Report(); got != "remove unused import" {
		t.Fatalf("unexpected self report message: %q", got)
	}
}

func TestSelfTraverseCovered(t *testing.T) {
	tr := Self.Traverse()
	if _, ok := tr["*ast.File"]; !ok {
		t.Fatal("expected self.Traverse to expose *ast.File visitor")
	}
}

func TestCollectImportsSkipsNonImportSpec(t *testing.T) {
	file := &ast.File{
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`},
					},
					// a non-ImportSpec in an import decl must be skipped
					&ast.ValueSpec{},
				},
			},
		},
	}
	imports := collectImports(file)
	if len(imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(imports))
	}
	if imports[0].path != `"fmt"` {
		t.Fatalf("unexpected path: %q", imports[0].path)
	}
}

// TestSelfFixRemovesUnusedImport covers self.Fix, the method the runner
// reflects on, by removing an unused import from a hand-built AST.
func TestSelfFixRemovesUnusedImport(t *testing.T) {
	file := &ast.File{
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			},
		},
	}
	Self.Fix(file, []types.Place{{Message: `remove unused import: "fmt"`}})
	if len(file.Decls) != 0 {
		t.Fatalf("expected import decl to be removed, got %d decls", len(file.Decls))
	}
}
