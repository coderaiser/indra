package convert_for_create_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_for_create"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-for-to-create-test", convert_for_create.Plugin{})

func TestConvertForToCreateTest(t *testing.T) {
	Test(t, "convert-for-to-create-test: report indratest.For call", func(t *T) {
		t.Report("convert-for-to-create-test", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform indratest.For to CreateTest", func(t *T) {
		t.Transform("convert-for-to-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report for CreateTest", func(t *T) {
		t.NoReport("already-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no transform for CreateTest", func(t *T) {
		t.NoTransform("already-create-test")
		t.End()
	})
}
