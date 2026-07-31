package compare

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"
)

func TestSentinels(t *testing.T) {
	a := ArgSlice{Args: []ast.Expr{ast.NewIdent("x")}}
	if a.Pos() != token.NoPos || a.End() != token.NoPos {
		t.Fatal("ArgSlice pos/end should be NoPos")
	}
	b := BodySlice{Stmts: []ast.Stmt{}}
	if b.Pos() != token.NoPos || b.End() != token.NoPos {
		t.Fatal("BodySlice pos/end should be NoPos")
	}

	// ArgSlice and BodySlice implement ast.Node so they fit in Vars.
	var _ ast.Node = a
	var _ ast.Node = b
}

func TestCompareInvalidPattern(t *testing.T) {
	if Compare(parseStmt(t, "foo()"), "((( invalid") != nil {
		t.Fatal("invalid pattern should not match")
	}
}

func TestMatchNodeNil(t *testing.T) {
	vars := make(Vars)
	if !matchNode(nil, nil, vars) {
		t.Fatal("nil pattern and nil real should match")
	}
	node := parseStmt(t, "foo()")
	if matchNode(nil, node, vars) {
		t.Fatal("nil pattern with non-nil real should not match")
	}
	if matchNode(node, nil, vars) {
		t.Fatal("non-nil pattern with nil real should not match")
	}
}

func TestMatchNodeIdentMismatch(t *testing.T) {
	vars := make(Vars)
	if matchNode(ast.NewIdent("x"), parseStmt(t, "foo()"), vars) {
		t.Fatal("ident pattern against non-ident real should not match")
	}
	if !matchNode(ast.NewIdent("x"), ast.NewIdent("x"), vars) {
		t.Fatal("equal idents should match")
	}
	if matchNode(ast.NewIdent("x"), ast.NewIdent("y"), vars) {
		t.Fatal("different idents should not match")
	}
}

func TestMatchChildrenNonStruct(t *testing.T) {
	vars := make(Vars)
	// nil pat: reflect Kind is Invalid -> return true
	if !matchChildren(nil, parseStmt(t, "foo()"), vars) {
		t.Fatal("non-struct pat should return true")
	}
	// real not a pointer struct -> return false
	if matchChildren(parseStmt(t, "foo()"), ArgSlice{}, vars) {
		t.Fatal("non-pointer real should return false")
	}
}

func TestPrintedNil(t *testing.T) {
	if printed(nil) != "" {
		t.Fatal("printed(nil) should be empty")
	}
	node := parseStmt(t, "foo()")
	if printed(node) == "" {
		t.Fatal("printed(node) should be non-empty")
	}
}

func TestMatchSliceLikeArguments(t *testing.T) {
	// non-ast.Expr slice elements fall through to the scalar comparison
	node := parseStmt(t, "foo()")
	call := node.(*ast.ExprStmt).X.(*ast.CallExpr)
	// CallExpr.Args is a slice of ast.Expr (interface elements).
	const tooMany = false
	_ = tooMany
	// mismatched element counts
	if Compare(node, "foo(1)") != nil {
		t.Fatal("arity mismatch should not match")
	}
	// empty arg lists (nil slices match as nil)
	if Compare(node, "foo()") == nil {
		t.Fatal("empty call should match")
	}
	_ = call
}

func TestMatchFieldDirect(t *testing.T) {
	vars := make(Vars)

	// slice both nil
	nilSlice := reflect.ValueOf(&struct{ S []ast.Expr }{S: nil}).Elem().Field(0)
	if !matchField("S", nilSlice, nilSlice, vars) {
		t.Fatal("nil slice should match")
	}

	// pointer nil accept
	nilPtr := reflect.ValueOf(&struct{ P *ast.Ident }{P: nil}).Elem().Field(0)
	if !matchField("P", nilPtr, nilPtr, vars) {
		t.Fatal("nil ptr should match nil ptr")
	}
	ptr := reflect.ValueOf(&struct{ P *ast.Ident }{P: ast.NewIdent("x")}).Elem().Field(0)
	if matchField("P", nilPtr, ptr, vars) {
		t.Fatal("nil ptr pattern must not match non-nil real")
	}
	if matchField("P", ptr, nilPtr, vars) {
		t.Fatal("non-nil ptr pattern must not match nil real")
	}

	// interface nil accept
	nilIf := reflect.ValueOf(&struct{ X ast.Expr }{X: nil}).Elem().Field(0)
	if !matchField("X", nilIf, nilIf, vars) {
		t.Fatal("nil interface should match nil interface")
	}
	if matchField("X", nilIf, ptr, vars) {
		t.Fatal("nil interface pattern must not match non-nil real")
	}
	if matchField("X", ptr, nilIf, vars) {
		t.Fatal("non-nil interface pattern must not match nil real")
	}

	// interface holding a non-ast.Node value
	anyA := reflect.ValueOf(&struct{ X any }{X: 5}).Elem().Field(0)
	anyB := reflect.ValueOf(&struct{ X any }{X: 5}).Elem().Field(0)
	anymis := reflect.ValueOf(&struct{ X any }{X: 6}).Elem().Field(0)
	if !matchField("X", anyA, anyB, vars) {
		t.Fatal("equal non-node interfaces should match")
	}
	if matchField("X", anyA, anymis, vars) {
		t.Fatal("different non-node interfaces should not match")
	}

	// default scalar comparison via token.Token field
	kindA := reflect.ValueOf(&struct{ K token.Token }{K: token.INT}).Elem().Field(0)
	kindB := reflect.ValueOf(&struct{ K token.Token }{K: token.INT}).Elem().Field(0)
	if !matchField("K", kindA, kindB, vars) {
		t.Fatal("equal scalar should match")
	}
}

func TestMatchSliceLikeScalarElement(t *testing.T) {
	// scalar (non-node) slice elements compare by value
	a := reflect.ValueOf(&struct{ S []int }{S: []int{1, 2}}).Elem().Field(0)
	b := reflect.ValueOf(&struct{ S []int }{S: []int{1, 2}}).Elem().Field(0)
	c := reflect.ValueOf(&struct{ S []int }{S: []int{1, 3}}).Elem().Field(0)
	vars := make(Vars)
	if !matchSliceLike(a, b, vars) {
		t.Fatal("equal int slices should match")
	}
	if matchSliceLike(a, c, vars) {
		t.Fatal("different int slices should not match")
	}
}
