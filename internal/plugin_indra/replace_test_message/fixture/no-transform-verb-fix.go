//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "no transform verb: some-fixture", func(t *T) {
		t.NoTransform("some-fixture")
		t.End()
	})
}
