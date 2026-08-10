package convert_for_to_create_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_for_to_create"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-for-to-create-test", convert_for_to_create.Plugin{})

func TestConvertForToCreateTest(t *testing.T) {
	Test(t, "convert-for-to-create-test: report: convert-for-to-create-test", func(t *T) {
		t.Report("convert-for-to-create-test", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform: convert-for-to-create-test", func(t *T) {
		t.Transform("convert-for-to-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report: already-create-test", func(t *T) {
		t.NoReport("already-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no transform: already-create-test", func(t *T) {
		t.NoTransform("already-create-test")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report: no-for-call", func(t *T) {
		t.NoReport("no-for-call")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no transform: no-for-call", func(t *T) {
		t.NoTransform("no-for-call")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report: non-ident-base", func(t *T) {
		t.NoReport("non-ident-base")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report: other-import", func(t *T) {
		t.NoReport("other-import")
		t.End()
	})

	Test(t, "convert-for-to-create-test: no report: other-method", func(t *T) {
		t.NoReport("other-method")
		t.End()
	})

	Test(t, "convert-for-to-create-test: report: non-t-star", func(t *T) {
		t.Report("non-t-star", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform: non-t-star", func(t *T) {
		t.Transform("non-t-star")
		t.End()
	})

	Test(t, "convert-for-to-create-test: report: has-t-star", func(t *T) {
		t.Report("has-t-star", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform: has-t-star", func(t *T) {
		t.Transform("has-t-star")
		t.End()
	})

	Test(t, "convert-for-to-create-test: report: mixed-calls", func(t *T) {
		t.Report("mixed-calls", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform: mixed-calls", func(t *T) {
		t.Transform("mixed-calls")
		t.End()
	})

	Test(t, "convert-for-to-create-test: report: other-import-path", func(t *T) {
		t.Report("other-import-path", "convert indratest.For to CreateTest")
		t.End()
	})

	Test(t, "convert-for-to-create-test: transform: other-import-path", func(t *T) {
		t.Transform("other-import-path")
		t.End()
	})
}
