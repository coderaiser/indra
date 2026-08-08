//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// non-literal-call: Test calls whose message is not a string literal or whose
// callback is not a function literal are skipped by hasMissingFixtureName.
var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, nonLiteralMessage, func(t *T) {
		t.End()
	})
	Test(t, "remove-skip: not-func", nonCallback)
}
