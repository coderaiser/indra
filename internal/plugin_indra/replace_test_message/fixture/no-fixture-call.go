//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: some message", func(t *T) {
		t.End()
	})
}
