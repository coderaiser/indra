package remove_only_test

import (
	"testing"

	loader "coderaiser/indra/engine_loader"
	"coderaiser/indra/internal/plugin_tape/remove_only"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-only", remove_only.Plugin{})

func TestRemoveOnly(t *testing.T) {
	Test(t, "remove-only: report: remove-only", func(t *T) {
		t.Report("remove-only", `Remove "Test.Only"`)
		t.End()
	})

	Test(t, "remove-only: transform: remove-only", func(t *T) {
		t.Transform("remove-only")
		t.End()
	})

	Test(t, "remove-only: no report: no-only", func(t *T) {
		t.NoReport("no-only")
		t.End()
	})

	Test(t, "remove-only: no transform: no-only", func(t *T) {
		t.NoTransform("no-only")
		t.End()
	})

	Test(t, "remove-only: no report: not-test", func(t *T) {
		t.NoReport("not-test")
		t.End()
	})

	Test(t, "remove-only: no transform: not-test", func(t *T) {
		t.NoTransform("not-test")
		t.End()
	})

	Test(t, "remove-only: report: with-options", func(t *T) {
		t.Report("with-options", `Remove "Test.Only"`)
		t.End()
	})

	Test(t, "remove-only: transform: with-options", func(t *T) {
		t.Transform("with-options")
		t.End()
	})

	Test(t, "remove-only: no report: allowed-receiver without options", func(t *T) {
		t.NoReport("allowed-receiver")
		t.End()
	})

	Test(t, "remove-only: no transform: allowed-receiver without options", func(t *T) {
		t.NoTransform("allowed-receiver")
		t.End()
	})
}

func TestRemoveOnlyAllowed(t *testing.T) {
	suite := CreateTestConfig("remove-only", remove_only.Plugin{}, loader.Config{
		"remove-only": {Enabled: true, Options: map[string]any{"allowed": []string{"Suite"}}},
	})

	suite(t, "remove-only: report: allowed-receiver", func(t *T) {
		t.Report("allowed-receiver", `Remove "Test.Only"`)
		t.End()
	})

	suite(t, "remove-only: transform: allowed-receiver", func(t *T) {
		t.Transform("allowed-receiver")
		t.End()
	})
}
