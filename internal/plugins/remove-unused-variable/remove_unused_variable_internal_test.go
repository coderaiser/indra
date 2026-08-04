package remove_unused_variable

import (
	"go/ast"
	"go/token"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestReportMessage(t *testing.T) {
	tape.Test(t, "report: nil node returns static message", func(t *tape.T) {
		t.Equal(Report(nil), "remove unused variable")
		t.End()
	})
}

func TestReportUsedVar(t *testing.T) {
	tape.Test(t, "report: used var falls through to static message", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.ExprStmt{X: ast.NewIdent("x")},
		}}
		t.Equal(Report(block), "remove unused variable")
		t.End()
	})
}
