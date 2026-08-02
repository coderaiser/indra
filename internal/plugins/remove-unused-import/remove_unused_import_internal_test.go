package remove_unused_import

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestPluginReportMessage(t *testing.T) {
	if got := Plugin.Report(); got != "remove unused import" {
		t.Fatalf("unexpected report message: %q", got)
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
