package compare

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// parseStmt parses a single statement from src and returns its node.
func parseStmt(t *testing.T, src string) ast.Node {
	t.Helper()
	wrap := "package p\nfunc _() {\n" + src + "\n}\n"
	file, err := parser.ParseFile(token.NewFileSet(), "x.go", wrap, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	fn := file.Decls[0].(*ast.FuncDecl)
	return fn.Body.List[0]
}

func named(t *testing.T, vars Vars, name string) ast.Node {
	t.Helper()
	if vars == nil {
		t.Fatal("expected non-nil match")
	}
	v, ok := vars[name]
	if !ok {
		t.Fatalf("expected hole %q to be bound", name)
	}
	return v
}

func TestCompareDiscard(t *testing.T) {
	node := parseStmt(t, "foo(1, 2)")
	vars := GetTemplateValues(node, "__")
	if vars == nil {
		t.Fatal("__ should match any expression")
	}
	if len(vars) != 0 {
		t.Fatal("__ must not store a binding")
	}
}

func TestCompareBindLinked(t *testing.T) {
	node := parseStmt(t, "t.Equal(x, x)")
	vars := GetTemplateValues(node, "t.Equal(__a, __a)")
	if vars == nil {
		t.Fatal("__a linked pattern should match")
	}
	a := named(t, vars, "__a")
	ident, ok := a.(*ast.Ident)
	if !ok || ident.Name != "x" {
		t.Fatalf("expected x, got %v", a)
	}

	// second occurrence of a different source must fail
	mismatch := parseStmt(t, "t.Equal(x, y)")
	if GetTemplateValues(mismatch, "t.Equal(__a, __a)") != nil {
		t.Fatal("linked hole with different values should not match")
	}
}

func TestCompareArgs(t *testing.T) {
	node := parseStmt(t, "t.Equal(a, b, c)")
	vars := GetTemplateValues(node, "t.Equal(__args)")
	if vars == nil {
		t.Fatal("__args pattern should match")
	}
	args, ok := named(t, vars, "__args").(ArgSlice)
	if !ok {
		t.Fatalf("expected ArgSlice, got %T", named(t, vars, "__args"))
	}
	if len(args.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args.Args))
	}
}

func TestCompareBody(t *testing.T) {
	src := `tape.Test(t, "x", func(t *tape.T) {
		t.Equal(1, 1)
		t.End()
	})`
	node := parseStmt(t, src)
	pattern := `tape.Test(__t, __, func(__t *tape.T) { __body })`
	vars := GetTemplateValues(node, pattern)
	if vars == nil {
		t.Fatal("__body pattern should match")
	}
	body, ok := named(t, vars, "__body").(BodySlice)
	if !ok {
		t.Fatalf("expected BodySlice, got %T", named(t, vars, "__body"))
	}
	if len(body.Stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(body.Stmts))
	}
}

func TestCompareArray(t *testing.T) {
	node := parseStmt(t, "t.DeepEqual(x, []int{1, 2})")
	vars := GetTemplateValues(node, "t.DeepEqual(__a, __array)")
	if vars == nil {
		t.Fatal("__array should match CompositeLit with ArrayType")
	}

	mapLit := parseStmt(t, "t.DeepEqual(x, map[string]int{})")
	if GetTemplateValues(mapLit, "t.DeepEqual(__a, __array)") != nil {
		t.Fatal("__array must reject a MapType composite literal")
	}

	nonLit := parseStmt(t, "t.DeepEqual(x, y)")
	if GetTemplateValues(nonLit, "t.DeepEqual(__a, __array)") != nil {
		t.Fatal("__array must reject a non-composite-literal argument")
	}
}

func TestCompareTypeMismatch(t *testing.T) {
	node := parseStmt(t, "foo(x)")
	if GetTemplateValues(node, "bar(x)") != nil {
		t.Fatal("different function names should not match")
	}
}

func TestCompareStructuralMismatch(t *testing.T) {
	node := parseStmt(t, "t.Equal(a, b)")
	if GetTemplateValues(node, "t.Equal(a)") != nil {
		t.Fatal("different arity should not match")
	}
}

func TestCompareEqual(t *testing.T) {
	node := parseStmt(t, "t.Equal(a, b)")
	vars := GetTemplateValues(node, "t.Equal(__a, __b)")
	if vars == nil {
		t.Fatal("t.Equal(__a, __b) should match")
	}
	if ai, ok := named(t, vars, "__a").(*ast.Ident); !ok || ai.Name != "a" {
		t.Fatal("expected __a to bind a")
	}
	if bi, ok := named(t, vars, "__b").(*ast.Ident); !ok || bi.Name != "b" {
		t.Fatal("expected __b to bind b")
	}
}

func TestCompareStruct(t *testing.T) {
	node := parseStmt(t, "t.DeepEqual(x, point{1, 2})")
	vars := GetTemplateValues(node, "t.DeepEqual(__a, __struct)")
	if vars == nil {
		t.Fatal("__struct should match CompositeLit with named type")
	}

	selLit := parseStmt(t, "t.DeepEqual(x, pkg.Point{})")
	if GetTemplateValues(selLit, "t.DeepEqual(__a, __struct)") == nil {
		t.Fatal("__struct should match a qualified struct composite literal")
	}

	sliceLit := parseStmt(t, "t.DeepEqual(x, []int{1})")
	if GetTemplateValues(sliceLit, "t.DeepEqual(__a, __struct)") != nil {
		t.Fatal("__struct must reject an ArrayType composite literal")
	}

	anonLit := parseStmt(t, "t.DeepEqual(x, struct{}{})")
	if GetTemplateValues(anonLit, "t.DeepEqual(__a, __struct)") != nil {
		t.Fatal("__struct must reject an anonymous struct composite literal")
	}

	nonLit := parseStmt(t, "t.DeepEqual(x, y)")
	if GetTemplateValues(nonLit, "t.DeepEqual(__a, __struct)") != nil {
		t.Fatal("__struct must reject a non-composite-literal argument")
	}
}
