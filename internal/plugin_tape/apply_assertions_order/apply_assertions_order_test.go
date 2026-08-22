package apply_assertions_order_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/apply_assertions_order"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("apply-assertions-order", apply_assertions_order.Plugin{})

func TestApplyAssertionsOrder(t *testing.T) {
	Test(t, "apply-assertions-order: report: apply-assertions-order", func(t *T) {
		t.Report("apply-assertions-order", "Apply assertions order")
		t.End()
	})

	Test(t, "apply-assertions-order: transform: apply-assertions-order", func(t *T) {
		t.Transform("apply-assertions-order")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: correct-order", func(t *T) {
		t.NoReport("correct-order")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: correct-order", func(t *T) {
		t.NoTransform("correct-order")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: single-assertion", func(t *T) {
		t.NoReport("single-assertion")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: single-assertion", func(t *T) {
		t.NoTransform("single-assertion")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: leading-gap", func(t *T) {
		t.NoReport("leading-gap")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: leading-gap", func(t *T) {
		t.NoTransform("leading-gap")
		t.End()
	})
	Test(t, "apply-assertions-order: no report: two-gaps", func(t *T) {
		t.NoReport("two-gaps")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: two-gaps", func(t *T) {
		t.NoTransform("two-gaps")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: index-recv", func(t *T) {
		t.NoReport("index-recv")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: index-recv", func(t *T) {
		t.NoTransform("index-recv")
		t.End()
	})
	Test(t, "apply-assertions-order: no report: end-first", func(t *T) {
		t.NoReport("end-first")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: end-first", func(t *T) {
		t.NoTransform("end-first")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: gap-only", func(t *T) {
		t.NoReport("gap-only")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: gap-only", func(t *T) {
		t.NoTransform("gap-only")
		t.End()
	})
	Test(t, "apply-assertions-order: no report: bare-expr", func(t *T) {
		t.NoReport("bare-expr")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: bare-expr", func(t *T) {
		t.NoTransform("bare-expr")
		t.End()
	})

	Test(t, "apply-assertions-order: no report: plain-call", func(t *T) {
		t.NoReport("plain-call")
		t.End()
	})

	Test(t, "apply-assertions-order: no transform: plain-call", func(t *T) {
		t.NoTransform("plain-call")
		t.End()
	})
}
