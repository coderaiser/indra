//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// no-separator: the message has no ": " rule prefix, so afterSeparator returns
// the whole string and the fixture name still gets appended.
var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "report missing", func(t *T) {
		t.Report("some-fix", "remove Test.Skip call")
		t.End()
	})
}
