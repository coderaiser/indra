package remove_unused_variables

import (
	"go/ast"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestReportOtherNode(t *testing.T) {
	Test(t, "report: non-file non-block node returns static message", func(t *T) {
		result := Report(ast.NewIdent("x"))
		
		t.Equal(result, "remove unused variable")
		t.End()
	})
}

func TestReportUnusedImport(t *testing.T) {
	Test(t, "report: importFinding returns unused import message", func(t *T) {
		file := unusedImportFile()
		imports := collectImports(file)
		finding := &importFinding{file: file, spec: imports[0].spec}
		result := Report(finding)
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
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}
		file := &ast.File{
			Name:    ast.NewIdent("fixture"),
			Imports: []*ast.ImportSpec{spec},
			Decls: []ast.Decl{
				&ast.GenDecl{
					Tok:   token.IMPORT,
					Specs: []ast.Spec{spec},
				},
			},
		}
		finding := &importFinding{file: file, spec: spec}
		Fix(finding, nil)
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

func TestImportFindingEnd(t *testing.T) {
	Test(t, "importFinding: End returns spec end position", func(t *T) {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}
		file := &ast.File{Name: ast.NewIdent("p"), Imports: []*ast.ImportSpec{spec}}
		finding := &importFinding{file: file, spec: spec}
		result := finding.End()
		t.Equal(result, spec.End())
		t.End()
	})
}

func TestUnusedVarNamesNonIdentLhs(t *testing.T) {
	Test(t, "unusedVarNames: non-ident lhs in := is skipped", func(t *T) {
		// e.g. a.b := 1 — not valid Go but AST can represent it
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{&ast.SelectorExpr{
					X:   ast.NewIdent("a"),
					Sel: ast.NewIdent("b"),
				}},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
		}}
		result := unusedVarNames(block)
		t.Equal(len(result), 0)
		t.End()
	})
}

func TestFixUnusedVarsNonValueSpec(t *testing.T) {
	Test(t, "fixUnusedVars: non-ValueSpec inside var GenDecl is kept", func(t *T) {
		// Construct a var GenDecl with a non-*ast.ValueSpec spec (defensive branch)
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok:   token.VAR,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			}},
		}}
		fixUnusedVars(block)
		// non-ValueSpec kept, AssignStmt with unused x dropped
		result := len(block.List)
		t.Equal(result, 1)
		t.End()
	})
}

func TestUnusedVarNamesNonGenDecl(t *testing.T) {
	Test(t, "unusedVarNames: non-GenDecl inside DeclStmt is skipped", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.DeclStmt{Decl: &ast.FuncDecl{Name: ast.NewIdent("f"), Type: &ast.FuncType{}, Body: &ast.BlockStmt{}}},
		}}
		result := unusedVarNames(block)
		t.Equal(len(result), 0)
		t.End()
	})
}

func TestUnusedVarNamesNonValueSpecInVar(t *testing.T) {
	Test(t, "unusedVarNames: non-ValueSpec inside var GenDecl is skipped", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok:   token.VAR,
				Specs: []ast.Spec{&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}}},
			}},
		}}
		result := unusedVarNames(block)
		t.Equal(len(result), 0)
		t.End()
	})
}

func TestFixUnusedVarsNonGenDecl(t *testing.T) {
	Test(t, "fixUnusedVars: non-GenDecl DeclStmt is kept", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.DeclStmt{Decl: &ast.FuncDecl{Name: ast.NewIdent("f"), Type: &ast.FuncType{}, Body: &ast.BlockStmt{}}},
		}}
		fixUnusedVars(block)
		result := len(block.List)
		t.Equal(result, 1)
		t.End()
	})
}

func TestUnusedVarNamesDuplicateDecl(t *testing.T) {
	Test(t, "unusedVarNames: duplicate var name in block counted once", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("x"), ast.NewIdent("x")}},
				},
			}},
		}}
		result := unusedVarNames(block)
		t.Equal(len(result), 1)
		t.End()
	})
}

func TestFixUnusedVarsPartialTupleUsed(t *testing.T) {
	Test(t, "fixUnusedVars: tuple with one used var keeps statement", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x"), ast.NewIdent("y")},
				Rhs: []ast.Expr{ast.NewIdent("f")},
			},
			&ast.ExprStmt{X: ast.NewIdent("y")},
		}}
		fixUnusedVars(block)
		result := len(block.List)
		t.Equal(result, 2)
		t.End()
	})
}

func TestFixUnusedVarsPartialVarSpecUsed(t *testing.T) {
	Test(t, "fixUnusedVars: var decl with one used spec keeps it", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("y")}},
				},
			}},
			&ast.ExprStmt{X: ast.NewIdent("y")},
		}}
		fixUnusedVars(block)
		// x removed, var y kept (used), ExprStmt kept
		result := len(block.List)
		t.Equal(result, 2)
		t.End()
	})
}

func TestUnusedVarNamesBlankVarName(t *testing.T) {
	Test(t, "unusedVarNames: blank var name is skipped", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("_")}},
				},
			}},
		}}
		result := unusedVarNames(block)
		t.Equal(len(result), 0)
		t.End()
	})
}
