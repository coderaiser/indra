package simplify_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/simplify"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/simplify", simplify.Plugin{})

func TestSimplify(t *testing.T) {
	Test(t, "conditions/simplify: report: report: simplify", func(t *T) {
		t.Report("simplify", "simplify condition")
		t.End()
	})

	Test(t, "conditions/simplify: transform: transform: simplify", func(t *T) {
		t.Transform("simplify")
		t.End()
	})

	Test(t, "conditions/simplify: no report: no report: no report: no-simplify-different", func(t *T) {
		t.NoReport("no-simplify-different")
		t.End()
	})

	Test(t, "conditions/simplify: no report: no report: no report: no-simplify-no-else", func(t *T) {
		t.NoReport("no-simplify-no-else")
		t.End()
	})

	Test(t, "conditions/simplify: no report: no report: no report: no-simplify-else-if", func(t *T) {
		t.NoReport("no-simplify-else-if")
		t.End()
	})
}
