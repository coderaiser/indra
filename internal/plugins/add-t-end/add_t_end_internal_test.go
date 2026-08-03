package add_t_end

import (
	"go/ast"
	"testing"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

func TestStmtsContainEnd(t *testing.T) {
	tape.Test(t, "stmtsContainEnd: nil list returns false", func(t *tape.T) {
		t.Equal(stmtsContainEnd(nil), false)
		t.End()
	})

	tape.Test(t, "stmtsContainEnd: End selector returns true", func(t *tape.T) {
		stmts := []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("End")},
		}}}
		t.Equal(stmtsContainEnd(stmts), true)
		t.End()
	})

	tape.Test(t, "stmtsContainEnd: AssignStmt returns false", func(t *tape.T) {
		t.Equal(stmtsContainEnd([]ast.Stmt{&ast.AssignStmt{}}), false)
		t.End()
	})

	tape.Test(t, "stmtsContainEnd: ExprStmt with non-CallExpr returns false", func(t *tape.T) {
		t.Equal(stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("foo")}}), false)
		t.End()
	})

	tape.Test(t, "stmtsContainEnd: CallExpr with non-SelectorExpr returns false", func(t *tape.T) {
		t.Equal(stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("foo")}}}), false)
		t.End()
	})

	tape.Test(t, "stmtsContainEnd: Equal selector returns false", func(t *tape.T) {
		stmts := []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("Equal")},
		}}}
		t.Equal(stmtsContainEnd(stmts), false)
		t.End()
	})
}

func TestGuard(t *testing.T) {
	var guard MatchFn
	for _, g := range Match() {
		guard = g
		break
	}

	tape.Test(t, "guard: non-BodySlice __body returns false", func(t *tape.T) {
		t.Equal(guard(Vars{"__body": ast.NewIdent("x")}), false)
		t.End()
	})

	tape.Test(t, "guard: empty BodySlice returns true", func(t *tape.T) {
		t.Equal(guard(Vars{"__body": compare.BodySlice{Stmts: []ast.Stmt{}}}), true)
		t.End()
	})
}
