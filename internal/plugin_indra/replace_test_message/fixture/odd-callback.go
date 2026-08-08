//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// odd-callback: extractFixtureName discriminators — a bare call, a fixture
// method with no arguments, and a fixture method with a non-literal argument.
var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: odd one", func(t *T) {
		helper()
		t.Transform()
		t.End()
	})
	Test(t, "remove-skip: odd two", func(t *T) {
		t.Report(someVar, "message")
		t.End()
	})
}
