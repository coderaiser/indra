package remove_skip_test

import (
	"testing"

	loader "coderaiser/indra/engine_loader"
	"coderaiser/indra/internal/plugin_tape/remove_skip"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-skip", remove_skip.Plugin{})

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report: remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: remove-skip", func(t *T) {
		t.Transform("remove-skip")
		t.End()
	})

	Test(t, "remove-skip: no report: no-skip", func(t *T) {
		t.NoReport("no-skip")
		t.End()
	})

	Test(t, "remove-skip: no transform: no-skip", func(t *T) {
		t.NoTransform("no-skip")
		t.End()
	})

	Test(t, "remove-skip: report: skip-skip", func(t *T) {
		t.Report("skip-skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: skip-skip", func(t *T) {
		t.Transform("skip-skip")
		t.End()
	})

	Test(t, "remove-skip: report: with-options", func(t *T) {
		t.Report("with-options", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: with-options", func(t *T) {
		t.Transform("with-options")
		t.End()
	})

	Test(t, "remove-skip: no report: not-test", func(t *T) {
		t.NoReport("not-test")
		t.End()
	})

	Test(t, "remove-skip: no transform: not-test", func(t *T) {
		t.NoTransform("not-test")
		t.End()
	})

	Test(t, "remove-skip: no report: allowed-receiver", func(t *T) {
		t.NoReport("allowed-receiver")
		t.End()
	})

	Test(t, "remove-skip: no transform: allowed-receiver", func(t *T) {
		t.NoTransform("allowed-receiver")
		t.End()
	})
}

// TestRemoveSkipAllowed exercises the option path: with allowed = ["Suite"],
// Suite.Skip calls are reported and rewritten too.
func TestRemoveSkipAllowed(t *testing.T) {
	suite := CreateTestConfig("remove-skip", remove_skip.Plugin{}, loader.Config{
		"remove-skip": {Enabled: true, Options: map[string]any{"allowed": []string{"Suite"}}},
	})

	suite(t, "remove-skip: report: allowed-receiver", func(t *T) {
		t.Report("allowed-receiver", "remove Test.Skip call")
		t.End()
	})

	suite(t, "remove-skip: transform: allowed-receiver", func(t *T) {
		t.Transform("allowed-receiver")
		t.End()
	})
}
