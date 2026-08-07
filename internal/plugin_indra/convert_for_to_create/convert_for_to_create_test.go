package convert_for_to_create_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_for_to_create"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-for-to-create-test", convert_for_to_create.Plugin{})

func TestConvertForToCreateTest(t *testing.T) {
	Test(t, "convert-for-to-create-test: report indratest.For call convert-for-to-create-test", func(t *T) {
		t.Report("convert-for-to-create-test", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform indratest.For to CreateTest convert-for-to-create-test", func(t *T) {
		t.Transform("convert-for-to-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report for CreateTest already-create-test", func(t *T) {
		t.NoReport("already-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no transform for CreateTest already-create-test", func(t *T) {
		t.NoTransform("already-create-test")
		t.End()
	})
}
