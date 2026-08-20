//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("short-calls", nil)

func f(t *testing.T) {
	// Two-arg Test — no callback, extractFixtureNameFromTest sees len < 3.
	Test(t, "short-calls: no callback")
	// Callback with zero-arg t.Report() — extractFixtureName sees len < 1.
	Test(t, "missing prefix", func(t *T) {
		t.Report()
		t.End()
	})
}