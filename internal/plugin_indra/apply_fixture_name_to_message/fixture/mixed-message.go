//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

// mixed-message has one unfixed message (pushing the file) and one Test call
// whose message is not a string literal (skipped by applyPrefix).
var Test = CreateTest("mixed-message", nil)

func f(t *testing.T) {
	Test(t, "missing fixture name", func(t *T) {
		t.End()
	})
	Test(t, someMessage, func(t *T) {
		t.End()
	})
}
