package remove_unused_import

import (
	"go/ast"
	"go/token"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestReportMessage(t *testing.T) {
	tape.Test(t, "report: returns remove unused import", func(t *tape.T) {
		t.Equal(Report(nil), "remove unused import")
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
		file := &ast.File{Name: ast.NewIdent("fixture"), Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			},
		}}
		Fix(file, nil)
		t.Equal(len(file.Decls), 0)
		t.End()
	})
}

// usedImportFile builds a file importing "fmt" and calling fmt.Println.
func usedImportFile() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("fixture"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("f"),
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{X: ast.NewIdent("fmt"), Sel: ast.NewIdent("Println")},
					}},
				}},
			},
		},
	}
}

func TestReportNoUnused(t *testing.T) {
	tape.Test(t, "Report: used import falls through to static message", func(t *tape.T) {
		t.Equal(Report(usedImportFile()), "remove unused import")
		t.End()
	})
}

func TestFixKeepsUsedImport(t *testing.T) {
	tape.Test(t, "Fix: keeps used import", func(t *tape.T) {
		file := usedImportFile()
		Fix(file, nil)
		t.Equal(len(file.Decls), 2)
		t.End()
	})
}

// blankImportFile builds a file importing "_ \"x\"".
func blankImportFile() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("fixture"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Name: ast.NewIdent("_"), Path: &ast.BasicLit{Kind: token.STRING, Value: `"x"`}},
				},
			},
		},
	}
}

func TestReportBlankImport(t *testing.T) {
	tape.Test(t, "Report: blank import falls through to static message", func(t *tape.T) {
		t.Equal(Report(blankImportFile()), "remove unused import")
		t.End()
	})
}

func TestFixKeepsBlankImport(t *testing.T) {
	tape.Test(t, "Fix: blank import is kept", func(t *tape.T) {
		file := blankImportFile()
		Fix(file, nil)
		t.Equal(len(file.Decls), 1)
		t.End()
	})
}
