//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("prefix-only", nil)

func f(t *testing.T) {
	Test(t, "prefix-only: some-fixture", func(t *T) {
		t.Report("some-fixture", "msg")
		t.End()
	})
}
