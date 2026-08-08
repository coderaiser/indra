package test_test

import (
	"testing"

	remove_skip "coderaiser/indra/internal/plugin_tape/remove_skip"
	. "coderaiser/indra/internal/test"
	"coderaiser/indra/types"
)

var Test = CreateTest("remove-skip", remove_skip.Plugin{})

var Group = ForGroup("g", []types.Rule{
	{Name: "remove-skip", Plugin: remove_skip.Plugin{}},
})

func TestShim(t *testing.T) {
	Test(t, "remove-skip: internal/test: shim CreateTest transforms fixture skip", func(t *T) {
		t.Transform("skip")
		t.End()
	})

	Test(t, "remove-skip: internal/test: shim CreateTest reports fixture skip", func(t *T) {
		t.Report("skip", "remove Test.Skip call")
		t.End()
	})

	Group(t, "internal/test: shim ForGroup transforms fixture skip", func(t *T) {
		t.Transform("skip")
		t.End()
	})
}
