package extract_result_from_assertion_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/extract_result_from_assertion"
	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.For("extract-result-from-assertion", extract_result_from_assertion.Plugin{})

func TestExtractResultFromAssertion(t *testing.T) {
	Test(t, "extract-result-from-assertion: report inline call", func(t *indratest.T) {
		t.Report("extract-result-from-assertion", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform inline call", func(t *indratest.T) {
		t.Transform("extract-result-from-assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: report array literal", func(t *indratest.T) {
		t.Report("array-literal", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform array literal", func(t *indratest.T) {
		t.Transform("array-literal")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report when already extracted", func(t *indratest.T) {
		t.NoReport("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform when already extracted", func(t *indratest.T) {
		t.NoTransform("extracted")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no report when result already declared", func(t *indratest.T) {
		t.NoReport("result-declared")
		t.End()
	})

	Test(t, "extract-result-from-assertion: no transform when result already declared", func(t *indratest.T) {
		t.NoTransform("result-declared")
		t.End()
	})
}
