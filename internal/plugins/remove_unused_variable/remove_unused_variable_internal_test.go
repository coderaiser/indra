package remove_unused_variable

import (
	"go/ast"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestReportMessage(t *testing.T) {
	Test(t, "report: nil node returns static message", func(t *T) {
		result := Report(nil)
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestReportUsedVar(t *testing.T) {
	Test(t, "report: used var falls through to static message", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.ExprStmt{X: ast.NewIdent("x")},
		}}
		result := Report(block)
		t.Equal(result, "remove unused variable")

		t.End()
	})
}
