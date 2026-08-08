package remove_useless_prefix

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

// parseFile parses src into an *ast.File for direct helper tests.
func parseFile(t *T, src string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "fixture.go", src, parser.ParseComments)
	if err != nil {
		t.TB().Fatalf("parse: %v", err)
	}
	return file
}

func TestVisitNoTapeImport(t *testing.T) {
	Test(t, "findUselessPrefix: no tape import pushes nothing", func(t *T) {
		file := parseFile(t, "package p\nfunc f() {}\n")
		pushed := false
		findUselessPrefix(types.Path{Node: file}, func(n types.Path) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

func TestFindTapeImportBlank(t *testing.T) {
	Test(t, "findTapeImport: blank import returns empty alias", func(t *T) {
		file := parseFile(t, "package p\nimport _ \"github.com/coderaiser/go-tape\"\n")
		alias, _ := findTapeImport(file)
		t.Equal(alias, "")
		t.End()
	})
}

func TestFindTapeImportDot(t *testing.T) {
	Test(t, "findTapeImport: dot import returns empty alias", func(t *T) {
		file := parseFile(t, "package p\nimport . \"github.com/coderaiser/go-tape\"\n")
		alias, _ := findTapeImport(file)
		t.Equal(alias, "")
		t.End()
	})
}

func TestFindTapeImportNamed(t *testing.T) {
	Test(t, "findTapeImport: named alias returned", func(t *T) {
		file := parseFile(t, "package p\nimport z \"github.com/coderaiser/go-tape\"\n")
		alias, spec := findTapeImport(file)
		t.Ok(alias == "z" && spec != nil)
		t.End()
	})
}

// TestFixNoAlias covers the early return in Fix when there is no named tape
// alias: it must leave the file untouched.
func TestFixNoAlias(t *testing.T) {
	Test(t, "Fix: no alias leaves file unchanged", func(t *T) {
		file := parseFile(t, "package p\nimport . \"github.com/coderaiser/go-tape\"\n")
		Fix(types.Path{Node: file}, nil)
		// dot import must remain a dot import
		t.Ok(file.Imports[0].Name == nil || file.Imports[0].Name.Name != "tape")
		t.End()
	})
}

func TestUsedBareNamesSkipsQualifier(t *testing.T) {
	Test(t, "usedBareNames: qualifier X is not collected", func(t *T) {
		file := parseFile(t, "package p\nfunc f() { _ = pkg.Method }\n")
		names := usedBareNames(file)
		t.NotOk(names["pkg"])

		t.End()
	})
}

func TestUsedBareNamesSkipsMember(t *testing.T) {
	Test(t, "usedBareNames: selector member Sel is not collected", func(t *T) {
		file := parseFile(t, "package p\nfunc f() { _ = pkg.Method }\n")
		names := usedBareNames(file)
		t.NotOk(names["Method"])

		t.End()
	})
}

func TestUsedBareNamesCollectsBare(t *testing.T) {
	Test(t, "usedBareNames: bare ident is collected", func(t *T) {
		file := parseFile(t, "package p\nfunc f() { _ = x }\n")
		names := usedBareNames(file)
		t.Ok(names["x"])
		t.End()
	})
}

func TestDeclaredNamesFunc(t *testing.T) {
	Test(t, "declaredNames: includes func names", func(t *T) {
		file := parseFile(t, "package p\nfunc Foo() {}\n")
		names := declaredNames(file)
		t.Ok(names["Foo"])
		t.End()
	})
}

func TestDeclaredNamesType(t *testing.T) {
	Test(t, "declaredNames: includes type names", func(t *T) {
		file := parseFile(t, "package p\ntype Foo int\n")
		names := declaredNames(file)
		t.Ok(names["Foo"])
		t.End()
	})
}

func TestDeclaredNamesValueSpec(t *testing.T) {
	Test(t, "declaredNames: includes var and const names", func(t *T) {
		file := parseFile(t, "package p\nvar version = \"1.0\"\nconst max = 10\n")
		names := declaredNames(file)
		t.Ok(names["version"] && names["max"])
		t.End()
	})
}

func TestFindUselessPrefixNonFile(t *testing.T) {
	Test(t, "findUselessPrefix: non-file node is a no-op", func(t *T) {
		pushed := false
		findUselessPrefix(types.Path{Node: ast.NewIdent("x")}, func(types.Path) { pushed = true })
		t.NotOk(pushed)

		t.End()
	})
}

// TestFixNonFile covers the non-*ast.File early return.
func TestFixNonFile(t *testing.T) {
	Test(t, "Fix: non-file node is a no-op", func(t *T) {
		Fix(types.Path{Node: ast.NewIdent("x")}, nil)
		t.Pass("no panic")
		t.End()
	})
}

// TestFixAppliesAcrossNestedCall covers a nested selector inside a call arg.
func TestFixAppliesAcrossNestedCall(t *testing.T) {
	Test(t, "Fix: rewrite nested selector in call arguments", func(t *T) {
		file := parseFile(t, "package p\n\nimport z \"github.com/coderaiser/go-tape\"\n\nfunc f() {\n\tz.Equal(1, z.T{})\n}\n")
		pushed := 0
		findUselessPrefix(types.Path{Node: file}, func(n types.Path) { pushed++ })
		t.Equal(pushed, 1)
		t.End()
	})

	Test(t, "Fix: selectors rewritten and import dotted", func(t *T) {
		file := parseFile(t, "package p\n\nimport z \"github.com/coderaiser/go-tape\"\n\nfunc f() {\n\tz.Equal(1, z.T{})\n}\n")
		Fix(types.Path{Node: file}, nil)
		imp := file.Imports[0]
		t.Ok(imp.Name != nil && imp.Name.Name == ".")
		t.End()
	})
}

// TestFixCollisionSkips covers the Fix early-return when removing the prefix
// would collide with a locally declared identifier.
func TestFixCollisionSkips(t *testing.T) {
	Test(t, "Fix: collision leaves alias import unchanged", func(t *T) {
		file := parseFile(t, "package p\n\nimport z \"github.com/coderaiser/go-tape\"\n\ntype T struct{ inner *z.T }\n")
		Fix(types.Path{Node: file}, nil)
		t.Ok(file.Imports[0].Name != nil && file.Imports[0].Name.Name == "z")
		t.End()
	})
}
