package remove_unused_variable_test

import (
	"go/ast"
	"go/token"
	"runtime"
	"testing"

	"coderaiser/indra/internal/plugins/remove_unused_variable"
	indratest "coderaiser/indra/internal/test"

	tape "github.com/coderaiser/go-tape"
)

var Test = indratest.CreateTest(runtime.Caller(0))

// TestReportDirect covers Report branches that are not reachable through the
// fixture harness (nil node, and a block whose vars are all used).
func TestReportDirect(t *testing.T) {
	tape.Test(t, "report: nil node returns static message", func(t *tape.T) {
		result := remove_unused_variable.Report(nil)
		t.Equal(result, "remove unused variable")
		t.End()
	})

	tape.Test(t, "report: no unused var returns static message", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.ExprStmt{X: ast.NewIdent("x")},
		}}
		result := remove_unused_variable.Report(block)
		t.Equal(result, "remove unused variable")
		t.End()
	})
}

func TestRemoveUnusedVariable(t *testing.T) {
	Test(t, "remove-unused-variable: report unused var", func(t *indratest.T) {
		t.Report("remove-unused-variable", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variable: no report when all vars used", func(t *indratest.T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variable: no report for blank in tuple assign", func(t *indratest.T) {
		t.NoReport("tuple-blank")
		t.End()
	})

	Test(t, "remove-unused-variable: fix removes unused var", func(t *indratest.T) {
		t.Transform("remove-unused-variable")
		t.End()
	})

	Test(t, "remove-unused-variable: fix blank the unused var in tuple", func(t *indratest.T) {
		t.Transform("tuple-unused")
		t.End()
	})

	Test(t, "remove-unused-variable: fix drops tuple with blank and unused var", func(t *indratest.T) {
		t.Transform("tuple-blank-unused")
		t.End()
	})
}
