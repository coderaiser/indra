package test_test

import (
	"testing"

	remove_skip "coderaiser/indra/internal/plugin_tape/remove_skip"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-skip", remove_skip.Plugin{})

func TestShim(t *testing.T) {
	Test(t, "remove-skip: internal/test: shim CreateTest transforms fixture skip", func(t *T) {
		t.Transform("skip")
		t.End()
	})

	Test(t, "remove-skip: internal/test: shim CreateTest reports fixture skip", func(t *T) {
		t.Report("skip", "remove Test.Skip call")
		t.End()
	})

	CreateTest("remove-skip", remove_skip.Plugin{})(t, "internal/test: shim CreateTest transforms fixture skip", func(t *T) {
		t.Transform("skip")
		t.End()
	})
}
