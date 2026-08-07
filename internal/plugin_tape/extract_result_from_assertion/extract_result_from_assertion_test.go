package extract_result_from_assertion_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/extract_result_from_assertion"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("extract-result-from-assertion", extract_result_from_assertion.Plugin{})

func TestExtractResultFromAssertion(t *testing.T) {
	Test(t, "extract-result-from-assertion: report inline call extract-result-from-assertion", func(t *T) {
		t.Report("extract-result-from-assertion", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform inline call extract-result-from-assertion", func(t *T) {
		t.Transform("extract-result-from-assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: report array literal array-literal", func(t *T) {
		t.Report("array-literal", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform array literal array-literal", func(t *T) {
		t.Transform("array-literal")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report when already extracted", func(t *T) {
		t.NoReport("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform when already extracted", func(t *T) {
		t.NoTransform("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report when result already declared result-declared", func(t *T) {
		t.NoReport("result-declared")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform when result already declared result-declared", func(t *T) {
		t.NoTransform("result-declared")
		t.End()
	})
}
