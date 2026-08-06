package convert_for_create_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_for_create"
	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.For("convert-for-to-create-test", convert_for_create.Plugin{})

func TestConvertForToCreateTest(t *testing.T) {
	Test(t, "convert-for-to-create-test: report indratest.For call", func(t *indratest.T) {
		t.Report("convert-for-to-create-test", "convert indratest.For to createTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform indratest.For to createTest", func(t *indratest.T) {
		t.Transform("convert-for-to-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report for createTest", func(t *indratest.T) {
		t.NoReport("already-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no transform for createTest", func(t *indratest.T) {
		t.NoTransform("already-create-test")
		t.End()
	})
}
