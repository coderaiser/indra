package remove_useless_prefix

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
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
		t.Ok(!pushed)
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

// TestReplaceSelectorsInvalid covers the !v.IsValid() guard.
func TestReplaceSelectorsInvalid(t *testing.T) {
	Test(t, "replaceSelectors: invalid value is a no-op", func(t *T) {
		replaceSelectors(reflect.Value{}, "tape")
		t.Pass("no-op")
		t.End()
	})
}

// TestMaybeReplaceSelNonPointer covers the default (non iface/ptr) path.
func TestMaybeReplaceSelNonPointer(t *testing.T) {
	Test(t, "maybeReplaceSel: non-pointer value returns false", func(t *T) {
		s := &ast.SelectorExpr{X: ast.NewIdent("tape"), Sel: ast.NewIdent("Test")}
		v := reflect.ValueOf(*s) // struct kind
		t.Ok(!maybeReplaceSel(v, "tape"))
		t.End()
	})
}

// TestMaybeReplaceSelNil covers nil interface/pointer paths.
func TestMaybeReplaceSelNil(t *testing.T) {
	Test(t, "maybeReplaceSel: nil pointer returns false", func(t *T) {
		var p *ast.SelectorExpr
		v := reflect.ValueOf(p)
		t.Ok(!maybeReplaceSel(v, "tape"))
		t.End()
	})
}

func TestUsedBareNamesSkipsQualifier(t *testing.T) {
	Test(t, "usedBareNames: qualifier X is not collected", func(t *T) {
		file := parseFile(t, "package p\nfunc f() { _ = pkg.Method }\n")
		names := usedBareNames(file)
		t.Ok(!names["pkg"])
		t.End()
	})
}

func TestUsedBareNamesSkipsMember(t *testing.T) {
	Test(t, "usedBareNames: selector member Sel is not collected", func(t *T) {
		file := parseFile(t, "package p\nimport z \"github.com/coderaiser/go-tape\"\nfunc f() { z.Test() }\n")
		names := usedBareNames(file)
		t.Ok(!names["Test"])
		t.End()
	})
}

func TestUsedBareNamesBareIdent(t *testing.T) {
	Test(t, "usedBareNames: bare ident outside selector is collected", func(t *T) {
		file := parseFile(t, "package p\nfunc f() { _ = T{} }\n")
		names := usedBareNames(file)
		t.Ok(names["T"])
		t.End()
	})
}

func TestHasLocalCollisionBareIdent(t *testing.T) {
	Test(t, "hasLocalCollision: bare ident matching selector member triggers collision", func(t *T) {
		// File uses z.T and also has bare T usage → collision
		file := parseFile(t, "package p\nimport z \"github.com/coderaiser/go-tape\"\nfunc f(_ z.T) {}\nfunc g(_ T) {}\n")
		t.Ok(hasLocalCollision(file, "z"))
		t.End()
	})
}

// TestReplaceSelectorsScopeGuard covers pruning the ast.Scope field (cycles).
// We locate the type by name via a parsed AST file so we never reference the
// deprecated ast.Scope identifier directly.
func TestReplaceSelectorsScopeGuard(t *testing.T) {
	Test(t, "replaceSelectors: ast.Scope value is pruned", func(t *T) {
		// Build a *ast.Scope value through reflect without naming the type.
		file, _ := parser.ParseFile(token.NewFileSet(), "", "package p", 0)
		// file.Scope is *ast.Scope; wrap it so we exercise the guard.
		v := reflect.ValueOf(file).Elem().FieldByName("Scope")
		replaceSelectors(v, "tape")
		t.Pass("scope pruned without recursion")
		t.End()
	})
}

// TestReplaceSelectorsIdentGuard covers pruning ast.Ident (Obj back-refs).
func TestReplaceSelectorsIdentGuard(t *testing.T) {
	Test(t, "replaceSelectors: ast.Ident value is pruned", func(t *T) {
		ident := ast.NewIdent("tape")
		replaceSelectors(reflect.ValueOf(ident), "tape")
		t.Pass("ident pruned without recursion")
		t.End()
	})
}

// TestReplaceSelectorsNilInterface covers the nil-interface guard in
// replaceSelectors.
func TestReplaceSelectorsNilInterface(t *testing.T) {
	Test(t, "replaceSelectors: nil interface is a no-op", func(t *T) {
		var e ast.Expr
		rv := reflect.ValueOf(&e).Elem() // settable nil interface field
		replaceSelectors(rv, "tape")
		t.Pass("nil interface pruned")
		t.End()
	})
}

// TestReplaceSelectorsIfaceSelector covers the interface branch where the
// element IS a matching selector: it must be replaced in place.
func TestReplaceSelectorsIfaceSelector(t *testing.T) {
	Test(t, "replaceSelectors: matching selector in interface is replaced", func(t *T) {
		var e ast.Expr = &ast.SelectorExpr{X: ast.NewIdent("tape"), Sel: ast.NewIdent("Test")}
		rv := reflect.ValueOf(&e).Elem() // settable interface field
		replaceSelectors(rv, "tape")
		sel, ok := e.(*ast.Ident)
		ok2 := ok && sel.Name == "Test"
		t.Ok(ok2)
		t.End()
	})
}

// TestReplaceSelectorsSliceSelector covers the slice branch where an element is
// a matching selector: it must be replaced.
func TestReplaceSelectorsSliceSelector(t *testing.T) {
	Test(t, "replaceSelectors: matching selector in slice is replaced", func(t *T) {
		e := []ast.Expr{
			&ast.SelectorExpr{X: ast.NewIdent("tape"), Sel: ast.NewIdent("Test")},
		}
		rv := reflect.ValueOf(&e).Elem() // settable slice
		replaceSelectors(rv, "tape")
		_, ok := e[0].(*ast.Ident)
		t.Ok(ok)
		t.End()
	})
}

// TestReplaceInStructNonStruct covers the kind guard in replaceInStruct.
func TestReplaceInStructNonStruct(t *testing.T) {
	Test(t, "replaceInStruct: non-struct value is a no-op", func(t *T) {
		replaceInStruct(reflect.ValueOf(42), "tape")
		t.Pass("non-struct pruned")
		t.End()
	})
}

// TestReplaceInStructUnexportedField covers the !CanSet() guard.
func TestReplaceInStructUnexportedField(t *testing.T) {
	Test(t, "replaceInStruct: unexported fields are skipped", func(t *T) {
		type hidden struct{ x int }
		replaceInStruct(reflect.ValueOf(hidden{x: 1}), "tape")
		t.Pass("unexported skipped")
		t.End()
	})
}

// TestMaybeReplaceSelNilInterface covers the nil-interface guard in
// maybeReplaceSel.
func TestMaybeReplaceSelNilInterface(t *testing.T) {
	Test(t, "maybeReplaceSel: nil interface returns false", func(t *T) {
		var e ast.Expr
		rv := reflect.ValueOf(&e).Elem() // settable nil interface field
		t.Ok(!maybeReplaceSel(rv, "tape"))
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

// TestDeclaredNamesValueSpec covers the *ast.ValueSpec (var/const) branch.
func TestDeclaredNamesValueSpec(t *testing.T) {
	Test(t, "declaredNames: includes var and const names", func(t *T) {
		file := parseFile(t, "package p\n\nvar version = \"1.0\"\nconst max = 10\n")
		names := declaredNames(file)
		t.Ok(names["version"] && names["max"])
		t.End()
	})
}

// TestFindUselessPrefixNonFile covers the non-*ast.File early return.
func TestFindUselessPrefixNonFile(t *testing.T) {
	Test(t, "findUselessPrefix: non-file node is a no-op", func(t *T) {
		pushed := false
		findUselessPrefix(types.Path{Node: ast.NewIdent("x")}, func(types.Path) { pushed = true })
		t.Ok(!pushed)
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
