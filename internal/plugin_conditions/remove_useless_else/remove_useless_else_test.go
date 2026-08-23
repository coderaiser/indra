package remove_useless_else_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/remove_useless_else"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/remove-useless-else", remove_useless_else.Plugin{})

func TestRemoveUselessElse(t *testing.T) {
	Test(t, "conditions/remove-useless-else: report: report: remove-useless-else", func(t *T) {
		t.Report("remove-useless-else", "remove useless else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: transform: transform: remove-useless-else", func(t *T) {
		t.Transform("remove-useless-else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: transform: transform: transform: break-else", func(t *T) {
		t.Transform("break-else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: transform: transform: transform: continue-else", func(t *T) {
		t.Transform("continue-else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: no report: no report: no report: no-useless-else", func(t *T) {
		t.NoReport("no-useless-else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: no report: no report: no report: no-return-else", func(t *T) {
		t.NoReport("no-return-else")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: transform: transform: transform: func-lit", func(t *T) {
		t.Transform("func-lit")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: no report: no report: no report: else-if", func(t *T) {
		t.NoReport("else-if")
		t.End()
	})

	Test(t, "conditions/remove-useless-else: no report: no report: no report: empty-body", func(t *T) {
		t.NoReport("empty-body")
		t.End()
	})
}
