package remove_unused_variables_test

import (
	"go/ast"
	"go/token"
	"testing"

	remove_unused_variables "coderaiser/indra/internal/plugin_remove_unused_variables"
	. "coderaiser/indra/internal/test"
	"coderaiser/indra/types"
)

var Test = CreateTest("remove-unused-variables", remove_unused_variables.Plugin{})

// TestReportDirect covers Report branches that are not reachable through the
// fixture harness (nil node, and a block whose vars are all used).
func TestReportDirect(t *testing.T) {
	Test(t, "remove-unused-variables: report: no unused var returns static message", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("x")},
				Rhs: []ast.Expr{ast.NewIdent("1")},
			},
			&ast.ExprStmt{X: ast.NewIdent("x")},
		}}
		result := remove_unused_variables.Report(types.Path{Node: block})

		t.Equal(result, "remove unused variable")
		t.End()
	})
}

func TestRemoveUnusedDeclarations(t *testing.T) {
	// import cases
	Test(t, "remove-unused-variables: report unused import remove-unused-import", func(t *T) {
		t.Report("remove-unused-import", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report when all imports used used-import", func(t *T) {
		t.NoReport("used-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for blank import blank-import", func(t *T) {
		t.NoReport("blank-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform for blank import blank-import", func(t *T) {
		t.NoTransform("blank-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for dot import dot-import", func(t *T) {
		t.NoReport("dot-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report unused aliased import alias-unused", func(t *T) {
		t.Report("alias-unused", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used alias alias-used", func(t *T) {
		t.NoReport("alias-used")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused import remove-unused-import", func(t *T) {
		t.Transform("remove-unused-import")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused aliased import alias-unused", func(t *T) {
		t.Transform("alias-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: fix keeps used import in mixed block mixed-imports", func(t *T) {
		t.Transform("mixed-imports")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used hyphenated path import hyphen-import", func(t *T) {
		t.NoReport("hyphen-import")
		t.End()
	})

	Test(t, "remove-unused-variables: fix keeps used hyphenated path import hyphen-import", func(t *T) {
		t.Transform("hyphen-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report when import package name differs from path basename plugin-named-package", func(t *T) {
		t.NoReport("plugin-named-package")
		t.End()
	})

	Test(t, "remove-unused-variables: fix keeps import when package name differs from path basename plugin-named-package", func(t *T) {
		t.NoTransform("plugin-named-package")
		t.End()
	})

	// go-prefix / version-suffix imports used under their real package name
	Test(t, "remove-unused-variables: no report go-prefix-import", func(t *T) {
		t.NoReport("go-prefix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform go-prefix-import", func(t *T) {
		t.NoTransform("go-prefix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report go-prefix-import-unused", func(t *T) {
		t.Report("go-prefix-import-unused", `remove unused import: "github.com/coderaiser/go-tape"`)
		t.End()
	})

	Test(t, "remove-unused-variables: transform go-prefix-import-unused", func(t *T) {
		t.Transform("go-prefix-import-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: no report version-suffix-import", func(t *T) {
		t.NoReport("version-suffix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform version-suffix-import", func(t *T) {
		t.NoTransform("version-suffix-import")
		t.End()
	})

	// multi go-prefix imports both used
	Test(t, "remove-unused-variables: no report multi-go-prefix", func(t *T) {
		t.NoReport("multi-go-prefix")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform multi-go-prefix", func(t *T) {
		t.NoTransform("multi-go-prefix")
		t.End()
	})

	// multi go-prefix, one unused
	Test(t, "remove-unused-variables: report multi-go-prefix-one-unused", func(t *T) {
		t.Report("multi-go-prefix-one-unused", `remove unused import: "github.com/coderaiser/go-coverage"`)
		t.End()
	})

	Test(t, "remove-unused-variables: transform multi-go-prefix-one-unused", func(t *T) {
		t.Transform("multi-go-prefix-one-unused")
		t.End()
	})

	// variable cases
	Test(t, "remove-unused-variables: report unused var remove-unused-variable", func(t *T) {
		t.Report("remove-unused-variable", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variables: no report when all vars used used-var", func(t *T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for blank in tuple assign tuple-blank", func(t *T) {
		t.NoReport("tuple-blank")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused var remove-unused-variable", func(t *T) {
		t.Transform("remove-unused-variable")
		t.End()
	})

	Test(t, "remove-unused-variables: fix blank the unused var in tuple tuple-unused", func(t *T) {
		t.Transform("tuple-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: fix drops tuple with blank and unused var tuple-blank-unused", func(t *T) {
		t.Transform("tuple-blank-unused")
		t.End()
	})

	// const cases (new)
	Test(t, "remove-unused-variables: report unused const unused-const", func(t *T) {
		t.Report("unused-const", "remove unused const: timeout")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used const used-const", func(t *T) {
		t.NoReport("used-const")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused const unused-const", func(t *T) {
		t.Transform("unused-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform for used const used-const", func(t *T) {
		t.NoTransform("used-const")
		t.End()
	})

	// mixed const declaration (some used, some unused)
	Test(t, "remove-unused-variables: report mixed const decl const-mixed", func(t *T) {
		t.Report("const-mixed", "remove unused const: timeout")
		t.End()
	})

	Test(t, "remove-unused-variables: transform keeps used const and removes unused const-mixed", func(t *T) {
		t.Transform("const-mixed")
		t.End()
	})

	// var declaration cases
	Test(t, "remove-unused-variables: report unused var decl var-decl", func(t *T) {
		t.Report("var-decl", "remove unused variable: a")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused var decl var-decl", func(t *T) {
		t.Transform("var-decl")
		t.End()
	})

	// mixed var declaration (some used, some unused)
	Test(t, "remove-unused-variables: report mixed var decl var-mixed", func(t *T) {
		t.Report("var-mixed", "remove unused variable: b")
		t.End()
	})

	Test(t, "remove-unused-variables: transform keeps used var and removes unused var-mixed", func(t *T) {
		t.Transform("var-mixed")
		t.End()
	})

	// blank var declaration (blank + unused)
	Test(t, "remove-unused-variables: report blank var decl blank-var-decl", func(t *T) {
		t.Report("blank-var-decl", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variables: transform removes blank var decl blank-var-decl", func(t *T) {
		t.Transform("blank-var-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: report blank var decl2 blank-var-decl2", func(t *T) {
		t.Report("blank-var-decl2", "remove unused variable: y")
		t.End()
	})

	Test(t, "remove-unused-variables: transform removes blank var decl2 blank-var-decl2", func(t *T) {
		t.Transform("blank-var-decl2")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for type decl type-decl", func(t *T) {
		t.NoReport("type-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform for type decl type-decl", func(t *T) {
		t.NoTransform("type-decl")
		t.End()
	})

	// multiple unused imports
	Test(t, "remove-unused-variables: report first of multiple unused imports multi-unused-import", func(t *T) {
		t.Report("multi-unused-import", `remove unused import: "encoding/json"`)
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes all unused imports multi-unused-import", func(t *T) {
		t.Transform("multi-unused-import")
		t.End()
	})

	// happy-path fixture (import + var + const all removed)
	Test(t, "remove-unused-variables: transform removes all unused variables, imports, and consts remove-unused-variables", func(t *T) {
		t.Transform("remove-unused-variables")
		t.End()
	})

	// unused private function cases
	Test(t, "remove-unused-variables: report unused private function unused-private-func", func(t *T) {
		t.Report("unused-private-func", "remove unused private function: unusedHelper")
		t.End()
	})

	Test(t, "remove-unused-variables: transform removes unused private function unused-private-func", func(t *T) {
		t.Transform("unused-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for called private function used-private-func", func(t *T) {
		t.NoReport("used-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform for called private function used-private-func", func(t *T) {
		t.NoTransform("used-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for private function used as value private-func-as-value", func(t *T) {
		t.NoReport("private-func-as-value")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for init function init-func", func(t *T) {
		t.NoReport("init-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for main function main-func", func(t *T) {
		t.NoReport("main-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for method method-func", func(t *T) {
		t.NoReport("method-func")
		t.End()
	})
}
