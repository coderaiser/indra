package extract_result_from_assertion_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/extract_result_from_assertion"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("extract-result-from-assertion", extract_result_from_assertion.Plugin{})

func TestExtractResultFromAssertion(t *testing.T) {
	Test(t, "extract-result-from-assertion: report: extract-result-from-assertion", func(t *T) {
		t.Report("extract-result-from-assertion", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform: extract-result-from-assertion", func(t *T) {
		t.Transform("extract-result-from-assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: report: array-literal", func(t *T) {
		t.Report("array-literal", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform: array-literal", func(t *T) {
		t.Transform("array-literal")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report: extracted", func(t *T) {
		t.NoReport("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform: extracted", func(t *T) {
		t.NoTransform("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report: result-declared", func(t *T) {
		t.NoReport("result-declared")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform: result-declared", func(t *T) {
		t.NoTransform("result-declared")
		t.End()
	})
	Test(t, "extract-result-from-assertion: no report: member-expr", func(t *T) {
		t.NoReport("member-expr")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform: member-expr", func(t *T) {
		t.NoTransform("member-expr")
		t.End()
	})

	Test(t, "extract-result-from-assertion: report: struct-literal", func(t *T) {
		t.Report("struct-literal", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform: struct-literal", func(t *T) {
		t.Transform("struct-literal")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report: both-declared", func(t *T) {
		t.NoReport("both-declared")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform: both-declared", func(t *T) {
		t.NoTransform("both-declared")
		t.End()
	})
}
