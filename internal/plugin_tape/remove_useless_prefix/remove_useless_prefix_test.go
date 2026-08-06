package remove_useless_prefix_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_useless_prefix"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-prefix", remove_useless_prefix.Plugin{})

func TestRemoveUselessPrefix(t *testing.T) {
	Test(t, "remove-useless-prefix: report named tape alias", func(t *T) {
		t.Report("remove-useless-prefix", "remove useless tape prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report for dot import", func(t *T) {
		t.NoReport("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report for no tape import", func(t *T) {
		t.NoReport("no-tape")
		t.End()
	})

	Test(t, "remove-useless-prefix: transform named to dot import", func(t *T) {
		t.Transform("remove-useless-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform for dot import", func(t *T) {
		t.NoTransform("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report when prefix collides with local decl", func(t *T) {
		t.NoReport("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform when prefix collides with local decl", func(t *T) {
		t.NoTransform("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report when member name used as bare ident", func(t *T) {
		t.NoReport("cross-file-t")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform when member name used as bare ident", func(t *T) {
		t.NoTransform("cross-file-t")
		t.End()
	})
}
