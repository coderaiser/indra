//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: transform: new-fixture", func(t *T) {
		t.Transform("new-fixture")
		t.End()
	})
}
