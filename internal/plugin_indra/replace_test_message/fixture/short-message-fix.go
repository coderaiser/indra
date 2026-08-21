//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// short-message: the first Test call has a verb mismatch (pushed by
// hasMismatch). The second Test call has a non-string BasicLit message
// (integer literal) whose RawText is shorter than 2 characters, so
// fixMessage returns early without modifying it.
var Test = CreateTest("short-message", nil)

func f(t *testing.T) {
	Test(t, "short-message: report: short-message", func(t *T) {
		t.Report("short-message")
	})
	Test(t, 0, func(t *T) {
		t.Report("short-message")
	})
}
