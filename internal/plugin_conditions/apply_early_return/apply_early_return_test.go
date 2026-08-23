package apply_early_return_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/apply_early_return"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/apply-early-return", apply_early_return.Plugin{})

func TestApplyEarlyReturn(t *testing.T) {
	Test(t, "conditions: apply-early-return: report", func(t *T) {
		t.Report("apply-early-return", "apply early return")
		t.End()
	})

	Test(t, "conditions: apply-early-return: transform", func(t *T) {
		t.Transform("apply-early-return")
		t.End()
	})

	Test(t, "conditions: apply-early-return: transform: multi-stmt", func(t *T) {
		t.Transform("multi-stmt")
		t.End()
	})

	Test(t, "conditions: apply-early-return: transform: var-func-lit", func(t *T) {
		t.Transform("var-func-lit")
		t.End()
	})

	Test(t, "conditions: apply-early-return: transform: func-lit", func(t *T) {
		t.Transform("func-lit")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: already-returns", func(t *T) {
		t.NoReport("already-returns")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: not-last", func(t *T) {
		t.NoReport("not-last")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: no-else", func(t *T) {
		t.NoReport("no-else")
		t.End()
	})

	Test(t, "conditions: apply-early-return: transform: empty-body", func(t *T) {
		t.Transform("empty-body")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: break-consequent", func(t *T) {
		t.NoReport("break-consequent")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: continue-consequent", func(t *T) {
		t.NoReport("continue-consequent")
		t.End()
	})

	Test(t, "conditions: apply-early-return: no report: else-if", func(t *T) {
		t.NoReport("else-if")
		t.End()
	})
}
