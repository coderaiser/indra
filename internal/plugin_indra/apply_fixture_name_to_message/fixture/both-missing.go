//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("both-missing", nil)

func f(t *testing.T) {
	Test(t, "msg", func(t *T) {
		t.Report("other-fixture")
	})
}