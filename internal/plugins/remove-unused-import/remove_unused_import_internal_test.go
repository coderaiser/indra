package remove_unused_import

import (
	"go/ast"
	"go/token"
	"testing"

	. "coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

func TestReportMessage(t *testing.T) {
	tape.Test(t, "report: returns remove unused import", func(t *tape.T) {
		t.Equal(Report(), "remove unused import")
		t.End()
	})
}

func TestCollectImportsLen(t *testing.T) {
	tape.Test(t, "collectImports: returns 1 entry for 1 ImportSpec", func(t *tape.T) {
		file := &ast.File{Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
					&ast.ValueSpec{},
				},
			},
		}}
		result := collectImports(file)
		t.Equal(len(result), 1)
		t.End()
	})
}

func TestCollectImportsPath(t *testing.T) {
	tape.Test(t, "collectImports: path equals import path literal", func(t *tape.T) {
		file := &ast.File{Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
		}}
		result := collectImports(file)
		t.Equal(result[0].path, `"fmt"`)
		t.End()
	})
}

func TestFixRemovesDecl(t *testing.T) {
	tape.Test(t, "Fix: removes GenDecl when all specs removed", func(t *tape.T) {
		file := &ast.File{Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			},
		}}
		Fix(file, []Place{{Message: `remove unused import: "fmt"`}})
		t.Equal(len(file.Decls), 0)
		t.End()
	})
}
