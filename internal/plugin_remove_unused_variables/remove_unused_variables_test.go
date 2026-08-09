package remove_unused_variables_test

import (
	"testing"

	remove_unused_variables "coderaiser/indra/internal/plugin_remove_unused_variables"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-unused-variables", remove_unused_variables.Plugin{})

func TestRemoveUnusedDeclarations(t *testing.T) {
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

	Test(t, "remove-unused-variables: no report go-prefix-import", func(t *T) {
		t.NoReport("go-prefix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report version-suffix-import", func(t *T) {
		t.NoReport("version-suffix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report unused const unused-const", func(t *T) {
		t.Report("unused-const", "remove unused const: timeout")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused const unused-const", func(t *T) {
		t.Transform("unused-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used const used-const", func(t *T) {
		t.NoReport("used-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no report for used var used-var", func(t *T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variables: report unused var decl var-decl", func(t *T) {
		t.Report("var-decl", "remove unused variable: a")
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes unused var decl var-decl", func(t *T) {
		t.Transform("var-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: report mixed var decl var-mixed", func(t *T) {
		t.Report("var-mixed", "remove unused variable: b")
		t.End()
	})

	Test(t, "remove-unused-variables: transform keeps used var and removes unused var-mixed", func(t *T) {
		t.Transform("var-mixed")
		t.End()
	})

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

	Test(t, "remove-unused-variables: report first of multiple unused imports multi-unused-import", func(t *T) {
		t.Report("multi-unused-import", `remove unused import: "encoding/json"`)
		t.End()
	})

	Test(t, "remove-unused-variables: fix removes all unused imports multi-unused-import", func(t *T) {
		t.Transform("multi-unused-import")
		t.End()
	})

	Test(t, "remove-unused-variables: transform removes all unused variables, imports, and consts remove-unused-variables", func(t *T) {
		t.Transform("remove-unused-variables")
		t.End()
	})

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
