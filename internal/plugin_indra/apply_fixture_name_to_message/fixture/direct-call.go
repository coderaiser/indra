//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("direct-call", nil)

func f(t *testing.T) {
	Test(t, "msg", func(t *T) {
		Report("some-fixture")
	})
}