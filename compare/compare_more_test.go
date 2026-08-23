package compare

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestMatchNodeBodySentinelBranches(t *testing.T) {
	vars := make(Vars)

	// pat block, single non-ExprStmt -> sentinel does not apply.
	retBlock := parseStmt(t, "{\nreturn\n}")
	if !matchNode(retBlock, parseStmt(t, "{\nreturn\n}"), vars) {
		t.Fatal("return-only blocks should match")
	}

	// pat block, single ExprStmt whose X is not a CallExpr.
	litBlock := parseStmt(t, "{\n1\n}")
	if !matchNode(litBlock, litBlock, vars) {
		t.Fatal("literal block should match")
	}

	// pat block, single call not named __body -> falls through to children.
	callBlock := parseStmt(t, "{\nfoo()\n}")
	if matchNode(callBlock, parseStmt(t, "{\nbar()\n}"), vars) {
		t.Fatal("different call blocks should not match")
	}

	// sentinel __body with non-block real must not capture.
	if matchNode(parseStmt(t, "{\n__body()\n}"), parseStmt(t, "foo()"), vars) {
		t.Fatal("__body sentinel with non-block real should not match")
	}
}

func TestMatchFieldInterfaceNilReal(t *testing.T) {
	vars := make(Vars)
	nonNil := reflect.ValueOf(&struct{ X ast.Expr }{X: ast.NewIdent("z")}).Elem().Field(0)
	nilIf := reflect.ValueOf(&struct{ X ast.Expr }{X: nil}).Elem().Field(0)
	if matchField("X", nonNil, nilIf, vars) {
		t.Fatal("non-nil interface pattern with nil real should not match")
	}
}

func TestMatchSliceLikeArgsHole(t *testing.T) {
	s := reflect.ValueOf(&struct{ S []ast.Expr }{S: []ast.Expr{ast.NewIdent("__args")}}).Elem().Field(0)
	dummy := reflect.ValueOf(&struct{ S []ast.Expr }{S: []ast.Expr{ast.NewIdent("x"), ast.NewIdent("y")}}).Elem().Field(0)
	vars := make(Vars)
	if !matchSliceLike(s, dummy, vars) {
		t.Fatal("a slice containing only __args should match any real slice")
	}
}

func TestMatchSliceLikeBlockStmts(t *testing.T) {
	pattern := parseStmt(t, "{\nfoo()\n}")
	if GetTemplateValues(parseStmt(t, "{\nfoo()\n}"), "{\nfoo()\n}") == nil {
		t.Fatal("block statement slices should match")
	}
	_ = pattern
}
