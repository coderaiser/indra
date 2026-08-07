//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "report Test.Skip call", func(t *T) {
		t.End()
	})
}
