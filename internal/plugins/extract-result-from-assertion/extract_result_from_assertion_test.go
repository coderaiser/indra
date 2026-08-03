package extract_result_from_assertion_test

import (
	"path/filepath"
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest("coderaiser/indra/internal/plugins/extract-result-from-assertion", _dir)
)

func TestExtractResultFromAssertion(t *testing.T) {
	Test(t, "extract-result-from-assertion: report inline call", func(t *indratest.T) {
		t.Report("inline-call", "extract inline expression from assertion")
		t.End()
	})

	Test(t, "extract-result-from-assertion: transform inline call", func(t *indratest.T) {
		t.Transform("inline-call")
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
}
