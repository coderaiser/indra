package add_t_end

import (
	"go/ast"
	"testing"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
	. "github.com/coderaiser/go-tape"
)

func TestStmtsContainEnd(t *testing.T) {
	Test(t, "stmtsContainEnd: nil list returns false", func(t *T) {
		result := stmtsContainEnd(nil)
		t.Equal(result, false)

		t.End()
	})

	Test(t, "stmtsContainEnd: End selector returns true", func(t *T) {
		stmts := []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("End")},
		}}}
		result := stmtsContainEnd(stmts)
		t.Equal(result, true)

		t.End()
	})

	Test(t, "stmtsContainEnd: AssignStmt returns false", func(t *T) {
		result := stmtsContainEnd([]ast.Stmt{&ast.AssignStmt{}})
		t.Equal(result, false)

		t.End()
	})

	Test(t, "stmtsContainEnd: ExprStmt with non-CallExpr returns false", func(t *T) {
		result := stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("foo")}})
		t.Equal(result, false)

		t.End()
	})

	Test(t, "stmtsContainEnd: CallExpr with non-SelectorExpr returns false", func(t *T) {
		result := stmtsContainEnd([]ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: ast.NewIdent("foo")}}})
		t.Equal(result, false)

		t.End()
	})

	Test(t, "stmtsContainEnd: Equal selector returns false", func(t *T) {
		stmts := []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("Equal")},
		}}}
		result := stmtsContainEnd(stmts)
		t.Equal(result, false)

		t.End()
	})
}

func TestGuard(t *testing.T) {
	var guard MatchFn
	for _, g := range Match() {
		guard = g
		break
	}

	Test(t, "guard: non-BodySlice __body returns false", func(t *T) {
		result := guard(Vars{"__body": ast.NewIdent("x")})
		t.Equal(result, false)

		t.End()
	})

	Test(t, "guard: empty BodySlice returns true", func(t *T) {
		result := guard(Vars{"__body": compare.BodySlice{Stmts: []ast.Stmt{}}})
		t.Equal(result, true)

		t.End()
	})
}
