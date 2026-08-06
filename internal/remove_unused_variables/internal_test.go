package remove_unused_variables

import (
	"go/ast"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
)

// ── Report ───────────────────────────────────────────────────────────────────

func TestReportNil(t *testing.T) {
	Test(t, "report: nil node returns static message", func(t *T) {
		result := Report(nil)
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestReportOtherNode(t *testing.T) {
	Test(t, "report: non-file non-block node returns static message", func(t *T) {
		result := Report(ast.NewIdent("x"))
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestReportUnusedImport(t *testing.T) {
	Test(t, "report: file returns unused import message", func(t *T) {
		result := Report(unusedImportFile())
		t.Equal(result, `remove unused import: "fmt"`)

		t.End()
	})
}

func TestReportNoUnused(t *testing.T) {
	Test(t, "report: clean file falls through to static message", func(t *T) {
		result := Report(usedImportFile())
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestReportUnusedConst(t *testing.T) {
	Test(t, "report: file with unused const reports const", func(t *T) {
		result := Report(unusedConstFile())
		t.Equal(result, "remove unused const: timeout")

		t.End()
	})
}

func TestReportUnusedVar(t *testing.T) {
	Test(t, "report: block with unused var reports variable", func(t *T) {
		result := Report(unusedVarBlock())
		t.Equal(result, "remove unused variable: x")

		t.End()
	})
}

func TestReportNoUnusedVar(t *testing.T) {
	Test(t, "report: block with all used vars returns static message", func(t *T) {
		result := Report(usedVarBlock())
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

// ── Traverse / traversal helpers ─────────────────────────────────────────────

func TestTraverseKeys(t *testing.T) {
	Test(t, "traverse: registers file and block visitors", func(t *T) {
		tr := Traverse()
		_, fileOK := tr["*ast.File"]
		_, blockOK := tr["*ast.BlockStmt"]
		t.Ok(fileOK && blockOK)

		t.End()
	})
}

func TestFindUnusedImportsPushes(t *testing.T) {
	Test(t, "findUnusedImportsAndConsts: pushes on unused import", func(t *T) {
		pushed := false
		findUnusedImportsAndConsts(unusedImportFile(), func(ast.Node) { pushed = true })
		t.Ok(pushed)

		t.End()
	})
}

func TestFindUnusedImportsNoPush(t *testing.T) {
	Test(t, "findUnusedImportsAndConsts: no push when clean", func(t *T) {
		pushed := false
		findUnusedImportsAndConsts(usedImportFile(), func(ast.Node) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

func TestFindUnusedConstsPushes(t *testing.T) {
	Test(t, "findUnusedImportsAndConsts: pushes on unused const", func(t *T) {
		pushed := false
		findUnusedImportsAndConsts(unusedConstFile(), func(ast.Node) { pushed = true })
		t.Ok(pushed)

		t.End()
	})
}

func TestFindUnusedVarsPushes(t *testing.T) {
	Test(t, "findUnusedVars: pushes on unused var", func(t *T) {
		pushed := false
		findUnusedVars(unusedVarBlock(), func(ast.Node) { pushed = true })
		t.Ok(pushed)

		t.End()
	})
}

func TestFindUnusedVarsNoPush(t *testing.T) {
	Test(t, "findUnusedVars: no push when vars used", func(t *T) {
		pushed := false
		findUnusedVars(usedVarBlock(), func(ast.Node) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

// ── collectImports ───────────────────────────────────────────────────────────

func TestCollectImportsLen(t *testing.T) {
	Test(t, "collectImports: returns 1 entry for 1 ImportSpec", func(t *T) {
		file := &ast.File{Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
					&ast.ValueSpec{},
				},
			},
		}}
		result := len(collectImports(file))
		t.Equal(result, 1)

		t.End()
	})
}

func TestCollectImportsPath(t *testing.T) {
	Test(t, "collectImports: path equals import path literal", func(t *T) {
		file := &ast.File{Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
		}}
		t.Equal(collectImports(file)[0].path, `"fmt"`)
		t.End()
	})
}

// ── Fix ──────────────────────────────────────────────────────────────────────

func TestFixRemovesDecl(t *testing.T) {
	Test(t, "fix: removes GenDecl when all specs removed", func(t *T) {
		file := &ast.File{Name: ast.NewIdent("fixture"), Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			},
		}}
		Fix(file, nil)
		result := len(file.Decls)
		t.Equal(result, 0)

		t.End()
	})
}

func TestFixKeepsUsedImport(t *testing.T) {
	Test(t, "fix: keeps used import", func(t *T) {
		file := usedImportFile()
		Fix(file, nil)
		result := len(file.Decls)
		t.Equal(result, 2)

		t.End()
	})
}

func TestFixKeepsBlankImport(t *testing.T) {
	Test(t, "fix: blank import is kept", func(t *T) {
		file := blankImportFile()
		Fix(file, nil)
		result := len(file.Decls)
		t.Equal(result, 1)

		t.End()
	})
}

func TestFixRemovesVars(t *testing.T) {
	Test(t, "fix: removes unused variables from block", func(t *T) {
		block := unusedVarBlock()
		Fix(block, nil)
		result := len(block.List)
		t.Equal(result, 0)

		t.End()
	})
}

func TestFixKeepsUsedVars(t *testing.T) {
	Test(t, "fix: keeps used variables in block", func(t *T) {
		block := usedVarBlock()
		Fix(block, nil)
		result := len(block.List)
		t.Equal(result, 2)

		t.End()
	})
}

func TestFixRemovesConsts(t *testing.T) {
	Test(t, "fix: removes unused const from file", func(t *T) {
		file := unusedConstFile()
		Fix(file, nil)
		result := len(file.Decls)
		t.Equal(result, 1)

		t.End()
	})
}

func TestPluginReport(t *testing.T) {
	Test(t, "plugin: Report delegates to package func", func(t *T) {
		result := Plugin{}.Report(nil)
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestPluginTraverse(t *testing.T) {
	Test(t, "plugin: Traverse registers two visitors", func(t *T) {
		result := len(Plugin{}.Traverse())
		t.Equal(result, 2)

		t.End()
	})
}

func TestPluginFix(t *testing.T) {
	Test(t, "plugin: Fix accepts nil", func(t *T) {
		Plugin{}.Fix(nil, nil)
		t.Pass("no panic")
		t.End()
	})
}

func TestReportBlankImportFallthrough(t *testing.T) {
	Test(t, "report: blank import skips to static message", func(t *T) {
		result := Report(blankImportFile())
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestCollectConstNamesSkipsNonValueSpec(t *testing.T) {
	Test(t, "collectConstNames: skips non-ValueSpec in const block", func(t *T) {
		file := &ast.File{Name: ast.NewIdent("fixture"), Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("timeout")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"x"`}},
				},
			},
		}}
		result := collectConstNames(file)
		t.Equal(len(result), 1)
		t.End()
	})
}

func TestFixConstsKeepsUsedAndSkipsNonValue(t *testing.T) {
	Test(t, "fixConsts: keeps used const and non-ValueSpec", func(t *T) {
		file := &ast.File{Name: ast.NewIdent("fixture"), Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("timeout")}},
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("limit")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"x"`}},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("f"),
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{List: []ast.Stmt{
					&ast.ExprStmt{X: ast.NewIdent("limit")},
				}},
			},
		}}
		fixUnusedConsts(file)
		genDecl := file.Decls[0].(*ast.GenDecl)
		result := len(genDecl.Specs)
		t.Equal(result, 2)

		t.End()
	})
}

// ── builders ─────────────────────────────────────────────────────────────────

func unusedImportFile() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("fixture"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
			&ast.FuncDecl{Name: ast.NewIdent("f"), Type: &ast.FuncType{}, Body: &ast.BlockStmt{}},
		},
	}
}

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

func unusedConstFile() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("fixture"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.CONST,
				Specs: []ast.Spec{&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("timeout")}}},
			},
			&ast.FuncDecl{Name: ast.NewIdent("f"), Type: &ast.FuncType{}, Body: &ast.BlockStmt{}},
		},
	}
}

func unusedVarBlock() *ast.BlockStmt {
	return &ast.BlockStmt{List: []ast.Stmt{
		&ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{ast.NewIdent("x")},
			Rhs: []ast.Expr{ast.NewIdent("1")},
		},
	}}
}

func usedVarBlock() *ast.BlockStmt {
	return &ast.BlockStmt{List: []ast.Stmt{
		&ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{ast.NewIdent("x")},
			Rhs: []ast.Expr{ast.NewIdent("1")},
		},
		&ast.ExprStmt{X: ast.NewIdent("x")},
	}}
}
