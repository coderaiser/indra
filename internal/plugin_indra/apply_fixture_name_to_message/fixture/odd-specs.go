//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

// odd-specs exercises the extractRuleName discriminator branches plus the
// non-literal Test-message path. It is never compiled (//go:build ignore), so
// the repeated package-level Test names are intentional.
var (
	// valid rule source, with a non-Test sibling name.
	Test, Other = CreateTest("odd-specs", nil), CreateTest("other", nil)
	// name-Test value that is not a call expression.
	X, Test = 5, 6
	// name-Test call whose function is not CreateTest.
	Y, Test = 7, foo()
	// name-Test CreateTest call whose first argument is not a string literal.
	Z, Test = 8, CreateTest(someArg, nil)
	// name-Test spec with more names than values.
	W, Test = 9
)

func f(t *testing.T) {
	Test(t, "odd-specs: covered", func(t *T) {
		t.End()
	})
	Test(t, someMessage, func(t *T) {
		t.End()
	})
}
