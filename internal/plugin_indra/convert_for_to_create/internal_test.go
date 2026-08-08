package convert_for_to_create

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

func parse(t *testing.T, src string) *ast.File {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "fixture.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

// ── Report / Traverse ────────────────────────────────────────────────────────

func TestReportString(t *testing.T) {
	Test(t, "report: returns conversion message", func(t *T) {
		result := Report(types.Path{})
		t.Equal(result, "convert indratest.For to CreateTest")

		t.End()
	})
}

func TestTraverseKey(t *testing.T) {
	Test(t, "traverse: registers file visitor", func(t *T) {
		_, ok := Traverse()["*ast.File"]
		t.Ok(ok)

		t.End()
	})
}

// ── findForCalls / hasForCall ───────────────────────────────────────────────

func TestFindForCallsNonFile(t *testing.T) {
	Test(t, "findForCalls: skips non-file node", func(t *T) {
		pushed := false
		findForCalls(types.Path{Node: ast.NewIdent("x")}, func(types.Path) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

func TestFindForCallsNoFor(t *testing.T) {
	Test(t, "findForCalls: no push without indratest.For", func(t *T) {
		pushed := false
		file := parse(t.TB(), "package fixture\nfunc f() {}\n")
		findForCalls(types.Path{Node: file}, func(types.Path) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

func TestFindForCallsPushes(t *testing.T) {
	Test(t, "findForCalls: pushes file with indratest.For", func(t *T) {
		pushed := false
		file := parse(t.TB(), "package fixture\nvar x = indratest.For(\"a\", f)\n")
		findForCalls(types.Path{Node: file}, func(types.Path) { pushed = true })
		t.Ok(pushed)

		t.End()
	})
}

func TestHasForCallTrue(t *testing.T) {
	Test(t, "hasForCall: detects indratest.For", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = indratest.For(\"a\", f)\n")
		t.Ok(hasForCall(types.Path{Node: file}))

		t.End()
	})
}

func TestHasForCallFalse(t *testing.T) {
	Test(t, "hasForCall: returns false without indratest.For", func(t *T) {
		file := parse(t.TB(), "package fixture\nfunc f() {}\n")
		t.NotOk(hasForCall(types.Path{Node: file}))

		t.End()
	})
}

func TestHasForCallNonIdentBase(t *testing.T) {
	Test(t, "hasForCall: skips non-ident selector base", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = a.b.For(\"a\")\n")
		t.NotOk(hasForCall(types.Path{Node: file}))

		t.End()
	})
}

// ── Fix / rewriteImport ──────────────────────────────────────────────────────

func TestFixNonFile(t *testing.T) {
	Test(t, "fix: skips non-file node", func(t *T) {
		Fix(types.Path{Node: ast.NewIdent("x")}, nil)
		t.Pass("no panic")
		t.End()
	})
}

func TestRewriteImportChanges(t *testing.T) {
	Test(t, "rewriteImport: turns alias into dot import", func(t *T) {
		file := parse(t.TB(), "package fixture\nimport indratest \"coderaiser/indra/internal/test\"\nfunc f() {}\n")
		rewriteImport(file)
		t.Equal(file.Imports[0].Name.Name, ".")
		t.End()
	})
}

func TestRewriteImportPathMismatch(t *testing.T) {
	Test(t, "rewriteImport: skips other import path", func(t *T) {
		file := parse(t.TB(), "package fixture\nimport indratest \"other/pkg\"\nfunc f() {}\n")
		rewriteImport(file)
		t.Equal(file.Imports[0].Name.Name, "indratest")
		t.End()
	})
}

func TestRewriteImportAliasMismatch(t *testing.T) {
	Test(t, "rewriteImport: skips non-indratest alias", func(t *T) {
		file := parse(t.TB(), "package fixture\nimport other \"coderaiser/indra/internal/test\"\nfunc f() {}\n")
		rewriteImport(file)
		t.Equal(file.Imports[0].Name.Name, "other")
		t.End()
	})
}

// ── rewriteCalls ─────────────────────────────────────────────────────────────

func TestRewriteCallsChanges(t *testing.T) {
	Test(t, "rewriteCalls: replaces indratest.For with CreateTest", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = indratest.For(\"a\", f)\n")
		rewriteCalls(types.Path{Node: file})
		result := printCallFun(file)
		t.Equal(result, "CreateTest")

		t.End()
	})
}

func TestRewriteCallsNotSelector(t *testing.T) {
	Test(t, "rewriteCalls: skips plain call", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = f()\n")
		rewriteCalls(types.Path{Node: file})
		result := printCallFun(file)
		t.Equal(result, "f")

		t.End()
	})
}

func TestRewriteCallsSelectorNotFor(t *testing.T) {
	Test(t, "rewriteCalls: skips other indratest method", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = indratest.G(\"a\")\n")
		rewriteCalls(types.Path{Node: file})
		result := printCallFun(file)
		t.Equal(result, "indratest.G")

		t.End()
	})
}

func TestRewriteCallsNonIdentSel(t *testing.T) {
	Test(t, "rewriteCalls: skips nested selector base", func(t *T) {
		file := parse(t.TB(), "package fixture\nvar x = a.b.For(\"a\")\n")
		rewriteCalls(types.Path{Node: file})
		result := printCallFun(file)
		t.Equal(result, "a.b.For")

		t.End()
	})
}

// ── rewriteT ─────────────────────────────────────────────────────────────────

func TestRewriteTChanges(t *testing.T) {
	Test(t, "rewriteT: replaces *indratest.T with *T", func(t *T) {
		file := parse(t.TB(), "package fixture\nfunc f(t *indratest.T) {}\n")
		rewriteT(types.Path{Node: file})
		result := printParamType(file)
		t.Equal(result, "*T")

		t.End()
	})
}

func TestRewriteTNotStar(t *testing.T) {
	Test(t, "rewriteT: skips non-star param", func(t *T) {
		file := parse(t.TB(), "package fixture\nfunc f(t indratest.T) {}\n")
		rewriteT(types.Path{Node: file})
		result := printParamType(file)
		t.Equal(result, "indratest.T")

		t.End()
	})
}

func TestRewriteTStarNonSelector(t *testing.T) {
	Test(t, "rewriteT: skips star over non-selector", func(t *T) {
		file := parse(t.TB(), "package fixture\nfunc f(p *int) {}\n")
		rewriteT(types.Path{Node: file})
		result := printParamType(file)
		t.Equal(result, "*int")

		t.End()
	})
}

func TestRewriteTNotT(t *testing.T) {
	Test(t, "rewriteT: skips other indratest type", func(t *T) {
		file := parse(t.TB(), "package fixture\nfunc f(t *indratest.U) {}\n")
		rewriteT(types.Path{Node: file})
		result := printParamType(file)
		t.Equal(result, "*indratest.U")

		t.End()
	})
}

// ── Plugin ───────────────────────────────────────────────────────────────────

func TestPluginReport(t *testing.T) {
	Test(t, "plugin: Report delegates", func(t *T) {
		result := Plugin{}.Report(types.Path{})
		t.Equal(result, "convert indratest.For to CreateTest")

		t.End()
	})
}

func TestPluginTraverse(t *testing.T) {
	Test(t, "plugin: Traverse registers file visitor", func(t *T) {
		result := len(Plugin{}.Traverse())
		t.Equal(result, 1)

		t.End()
	})
}

func printCallFun(file *ast.File) string {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				call, ok := val.(*ast.CallExpr)
				if ok {
					return printNode(call.Fun)
				}
			}
		}
	}
	return ""
}

func printParamType(file *ast.File) string {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
			continue
		}
		return printNode(fn.Type.Params.List[0].Type)
	}
	return ""
}

func printNode(n ast.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	_ = format.Node(&buf, token.NewFileSet(), n)
	return buf.String()
}
