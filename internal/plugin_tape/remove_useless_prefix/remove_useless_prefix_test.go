package remove_useless_prefix_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_useless_prefix"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-prefix", remove_useless_prefix.Plugin{})

func TestRemoveUselessPrefix(t *testing.T) {
	Test(t, "remove-useless-prefix: report: remove-useless-prefix", func(t *T) {
		t.Report("remove-useless-prefix", "remove useless tape prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: no-prefix", func(t *T) {
		t.NoReport("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: no-tape", func(t *T) {
		t.NoReport("no-tape")
		t.End()
	})

	Test(t, "remove-useless-prefix: transform: remove-useless-prefix", func(t *T) {
		t.Transform("remove-useless-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: no-prefix", func(t *T) {
		t.NoTransform("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: local-collision", func(t *T) {
		t.NoReport("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: local-collision", func(t *T) {
		t.NoTransform("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: cross-file-t", func(t *T) {
		t.NoReport("cross-file-t")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: cross-file-t", func(t *T) {
		t.NoTransform("cross-file-t")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: blank-import", func(t *T) {
		t.NoReport("blank-import")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: blank-import", func(t *T) {
		t.NoTransform("blank-import")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: dot-import", func(t *T) {
		t.NoReport("dot-import")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: dot-import", func(t *T) {
		t.NoTransform("dot-import")
		t.End()
	})

	Test(t, "remove-useless-prefix: report: nested-selector", func(t *T) {
		t.Report("nested-selector", "remove useless tape prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: transform: nested-selector", func(t *T) {
		t.Transform("nested-selector")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report: var-collision", func(t *T) {
		t.NoReport("var-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform: var-collision", func(t *T) {
		t.NoTransform("var-collision")
		t.End()
	})
}
