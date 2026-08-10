package remove_unused_variables_test

import (
	"testing"

	remove_unused_variables "coderaiser/indra/internal/plugin_remove_unused_variables"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-unused-variables", remove_unused_variables.Plugin{})

func TestRemoveUnusedDeclarations(t *testing.T) {
	Test(t, "remove-unused-variables: report: remove-unused-import", func(t *T) {
		t.Report("remove-unused-import", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report: used-import", func(t *T) {
		t.NoReport("used-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: blank-import", func(t *T) {
		t.NoReport("blank-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform: blank-import", func(t *T) {
		t.NoTransform("blank-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: dot-import", func(t *T) {
		t.NoReport("dot-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report: alias-unused", func(t *T) {
		t.Report("alias-unused", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-variables: no report: alias-used", func(t *T) {
		t.NoReport("alias-used")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: remove-unused-import", func(t *T) {
		t.Transform("remove-unused-import")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: alias-unused", func(t *T) {
		t.Transform("alias-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: mixed-imports", func(t *T) {
		t.Transform("mixed-imports")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: hyphen-import", func(t *T) {
		t.NoReport("hyphen-import")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: hyphen-import", func(t *T) {
		t.Transform("hyphen-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: plugin-named-package", func(t *T) {
		t.NoReport("plugin-named-package")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform: plugin-named-package", func(t *T) {
		t.NoTransform("plugin-named-package")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: go-prefix-import", func(t *T) {
		t.NoReport("go-prefix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: version-suffix-import", func(t *T) {
		t.NoReport("version-suffix-import")
		t.End()
	})

	Test(t, "remove-unused-variables: report: unused-const", func(t *T) {
		t.Report("unused-const", "remove unused const: timeout")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: unused-const", func(t *T) {
		t.Transform("unused-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: used-const", func(t *T) {
		t.NoReport("used-const")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: used-var", func(t *T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variables: report: var-decl", func(t *T) {
		t.Report("var-decl", "remove unused variable: a")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: var-decl", func(t *T) {
		t.Transform("var-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: report: var-mixed", func(t *T) {
		t.Report("var-mixed", "remove unused variable: b")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: var-mixed", func(t *T) {
		t.Transform("var-mixed")
		t.End()
	})

	Test(t, "remove-unused-variables: report: blank-var-decl", func(t *T) {
		t.Report("blank-var-decl", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: blank-var-decl", func(t *T) {
		t.Transform("blank-var-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: report: blank-var-decl2", func(t *T) {
		t.Report("blank-var-decl2", "remove unused variable: y")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: blank-var-decl2", func(t *T) {
		t.Transform("blank-var-decl2")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: type-decl", func(t *T) {
		t.NoReport("type-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform: type-decl", func(t *T) {
		t.NoTransform("type-decl")
		t.End()
	})

	Test(t, "remove-unused-variables: report: multi-unused-import", func(t *T) {
		t.Report("multi-unused-import", `remove unused import: "encoding/json"`)
		t.End()
	})

	Test(t, "remove-unused-variables: transform: multi-unused-import", func(t *T) {
		t.Transform("multi-unused-import")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: remove-unused-variables", func(t *T) {
		t.Transform("remove-unused-variables")
		t.End()
	})

	Test(t, "remove-unused-variables: report: unused-private-func", func(t *T) {
		t.Report("unused-private-func", "remove unused private function: unusedHelper")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: unused-private-func", func(t *T) {
		t.Transform("unused-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: used-private-func", func(t *T) {
		t.NoReport("used-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no transform: used-private-func", func(t *T) {
		t.NoTransform("used-private-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: private-func-as-value", func(t *T) {
		t.NoReport("private-func-as-value")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: init-func", func(t *T) {
		t.NoReport("init-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: main-func", func(t *T) {
		t.NoReport("main-func")
		t.End()
	})

	Test(t, "remove-unused-variables: no report: method-func", func(t *T) {
		t.NoReport("method-func")
		t.End()
	})

	Test(t, "remove-unused-variables: report: assign-partial-unused", func(t *T) {
		t.Report("assign-partial-unused", "remove unused variable: b")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: assign-partial-unused", func(t *T) {
		t.Transform("assign-partial-unused")
		t.End()
	})

	Test(t, "remove-unused-variables: report: assign-blank-unused", func(t *T) {
		t.Report("assign-blank-unused", "remove unused variable: b")
		t.End()
	})

	Test(t, "remove-unused-variables: transform: assign-blank-unused", func(t *T) {
		t.Transform("assign-blank-unused")
		t.End()
	})
}
