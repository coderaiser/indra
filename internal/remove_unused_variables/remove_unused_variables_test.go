package remove_unused_variables_test

import (
	"go/ast"
	"go/token"
	"testing"

	"coderaiser/indra/internal/remove_unused_variables"
	indratest "coderaiser/indra/internal/test"

	tape "github.com/coderaiser/go-tape"
)

var Test = indratest.For("remove-unused-variables", remove_unused_variables.Plugin{})

// TestReportDirect covers Report branches that are not reachable through the
// fixture harness (nil node, and a block whose vars are all used).
func TestReportDirect(t *testing.T) {
	tape.Test(t, "report: nil node returns static message", func(t *tape.T) {
		result := remove_unused_variables.Report(nil)
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
		result := remove_unused_variables.Report(block)
		t.Equal(result, "remove unused variable")

		t.End()
	})
}

func TestRemoveUnusedDeclarations(t *testing.T) {
	// import cases
	Test(t, "remove-unused-variables: report unused import", func(t *indratest.T) {
		t.Report("remove-unused-import", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report when all imports used", func(t *indratest.T) {
		t.NoReport("used-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for blank import", func(t *indratest.T) {
		t.NoReport("blank-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for dot import", func(t *indratest.T) {
		t.NoReport("dot-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report unused aliased import", func(t *indratest.T) {
		t.Report("alias-unused", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used alias", func(t *indratest.T) {
		t.NoReport("alias-used")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused import", func(t *indratest.T) {
		t.Transform("remove-unused-import")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused aliased import", func(t *indratest.T) {
		t.Transform("alias-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: fix keeps used import in mixed block", func(t *indratest.T) {
		t.Transform("mixed-imports")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used hyphenated path import", func(t *indratest.T) {
		t.NoReport("hyphen-import")
		t.End()
	})

	Test(t, "remove-unused-variables: fix keeps used hyphenated path import", func(t *indratest.T) {
		t.Transform("hyphen-import")
		t.End()
	})

	// variable cases
	Test(t, "remove-unused-variables: report unused var", func(t *indratest.T) {
		t.Report("remove-unused-variable", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variables: no report when all vars used", func(t *indratest.T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for blank in tuple assign", func(t *indratest.T) {
		t.NoReport("tuple-blank")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused var", func(t *indratest.T) {
		t.Transform("remove-unused-variable")
		t.End()
	})

	Test(t, "remove-unused-variables: fix blank the unused var in tuple", func(t *indratest.T) {
		t.Transform("tuple-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: fix drops tuple with blank and unused var", func(t *indratest.T) {
		t.Transform("tuple-blank-unused")
		t.End()
	})

	// const cases (new)
	Test(t, "remove-unused-variables: report unused const", func(t *indratest.T) {
		t.Report("unused-const", "remove unused const: timeout")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used const", func(t *indratest.T) {
		t.NoReport("used-const")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused const", func(t *indratest.T) {
		t.Transform("unused-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform for used const", func(t *indratest.T) {
		t.NoTransform("used-const")
		t.End()
	})

	// happy-path fixture (import + var + const all removed)
	Test(t, "remove-unused-variables: transform removes all unused variables, imports, and consts", func(t *indratest.T) {
		t.Transform("remove-unused-variables")
		t.End()
	})
}
