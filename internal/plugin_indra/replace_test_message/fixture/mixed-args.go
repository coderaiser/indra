//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// mixed-args: pushed by the first Test call; the later odd-shaped calls are
// skipped by applyFixtureNames.
var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) {
		t.Report("some-fix", "remove Test.Skip call")
		t.End()
	})
	Test(t, nonLiteralMessage, func(t *T) {
		t.End()
	})
	Test(t, "remove-skip: third", nonCallback)
	Test(t, "remove-skip: fourth", func(t *T) {
		t.End()
	})
}
