package add_t_end

import (
	"go/ast"
	"testing"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
)

func TestStmtsContainEnd(t *testing.T) {
	// empty statement list never contains End
	if stmtsContainEnd(nil) {
		t.Error("empty stmts should not contain End")
	}

	// a call whose selector name is End
	if !stmtsContainEnd([]ast.Stmt{
		&ast.ExprStmt{X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("End")},
		}},
	}) {
		t.Error("expected End to be found")
	}

	// non-ExprStmt statements are skipped
	if stmtsContainEnd([]ast.Stmt{&ast.AssignStmt{}}) {
		t.Error("AssignStmt is not an End call")
	}

	// ExprStmt whose X is not a CallExpr is skipped
	if stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("foo")}}) {
		t.Error("ident is not an End call")
	}

	// CallExpr whose Fun is not a SelectorExpr is skipped
	if stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("foo")}}}) {
		t.Error("plain call is not an End call")
	}

	// non-End selector call does not count
	if stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
		Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("Equal")},
	}}}) {
		t.Error("Equal call should not count as End")
	}
}

func TestGuardRejectsNonBodySlice(t *testing.T) {
	// the guard must reject when __body is not bound to a BodySlice
	var guard MatchFn
	for _, g := range Match() {
		guard = g
	}
	vars := Vars{"__body": ast.NewIdent("x")}
	if guard(vars) {
		t.Error("guard should reject a non-BodySlice __body")
	}

	// with a genuine BodySlice the guard inspects its statements;
	// an empty body (no End present) must be reported as matching
	body := compare.BodySlice{Stmts: []ast.Stmt{}}
	if !guard(Vars{"__body": body}) {
		t.Error("guard should match an empty body (no End present)")
	}
}
