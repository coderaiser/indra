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
}
