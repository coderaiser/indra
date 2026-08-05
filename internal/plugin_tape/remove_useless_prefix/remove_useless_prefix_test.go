package remove_useless_prefix_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestRemoveUselessPrefix(t *testing.T) {
	Test(t, "remove-useless-prefix: report named tape alias", func(t *indratest.T) {
		t.Report("remove-useless-prefix", "remove useless tape prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report for dot import", func(t *indratest.T) {
		t.NoReport("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report for no tape import", func(t *indratest.T) {
		t.NoReport("no-tape")
		t.End()
	})

	Test(t, "remove-useless-prefix: transform named to dot import", func(t *indratest.T) {
		t.Transform("remove-useless-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform for dot import", func(t *indratest.T) {
		t.NoTransform("no-prefix")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report when prefix collides with local decl", func(t *indratest.T) {
		t.NoReport("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform when prefix collides with local decl", func(t *indratest.T) {
		t.NoTransform("local-collision")
		t.End()
	})

	Test(t, "remove-useless-prefix: no report when member name used as bare ident", func(t *indratest.T) {
		t.NoReport("cross-file-t")
		t.End()
	})

	Test(t, "remove-useless-prefix: no transform when member name used as bare ident", func(t *indratest.T) {
		t.NoTransform("cross-file-t")
		t.End()
	})
}
