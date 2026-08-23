package remove_boolean_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/remove_boolean"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/remove-boolean", remove_boolean.Plugin{})

func TestRemoveBoolean(t *testing.T) {
	Test(t, "conditions/remove-boolean: report: report: remove-boolean", func(t *T) {
		t.Report("remove-boolean", "remove boolean")
		t.End()
	})

	Test(t, "conditions/remove-boolean: transform: transform: remove-boolean", func(t *T) {
		t.Transform("remove-boolean")
		t.End()
	})

	Test(t, "conditions/remove-boolean: transform: transform: transform: not-equal-false", func(t *T) {
		t.Transform("not-equal-false")
		t.End()
	})

	Test(t, "conditions/remove-boolean: transform: transform: transform: equal-false", func(t *T) {
		t.Transform("equal-false")
		t.End()
	})

	Test(t, "conditions/remove-boolean: transform: transform: transform: not-equal-true", func(t *T) {
		t.Transform("not-equal-true")
		t.End()
	})

	Test(t, "conditions/remove-boolean: no report: no report: no report: no-boolean", func(t *T) {
		t.NoReport("no-boolean")
		t.End()
	})

	Test(t, "conditions/remove-boolean: no report: no report: no report: arithmetic", func(t *T) {
		t.NoReport("arithmetic")
		t.End()
	})

	Test(t, "conditions/remove-boolean: no report: no report: no report: not-boolean-literal", func(t *T) {
		t.NoReport("not-boolean-literal")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: expr-operand", func(t *T) {
		t.Transform("expr-operand")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: unary-operand", func(t *T) {
		t.Transform("unary-operand")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: func-lit-param", func(t *T) {
		t.Transform("func-lit-param")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: typed-var", func(t *T) {
		t.Transform("typed-var")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: literal-var", func(t *T) {
		t.Transform("literal-var")
		t.End()
	})

	Test(t, "conditions: remove-boolean: no report: interface-value", func(t *T) {
		t.NoReport("interface-value")
		t.End()
	})

	Test(t, "conditions: remove-boolean: no report: sum", func(t *T) {
		t.NoReport("sum")
		t.End()
	})

	Test(t, "conditions: remove-boolean: no report: unbound", func(t *T) {
		t.NoReport("unbound")
		t.End()
	})

	Test(t, "conditions: remove-boolean: transform: short-literal", func(t *T) {
		t.Transform("short-literal")
		t.End()
	})

	Test(t, "conditions: remove-boolean: no report: local-type", func(t *T) {
		t.NoReport("local-type")
		t.End()
	})
}
