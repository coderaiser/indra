//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-t-end", nil)

func f(t *testing.T) {
	Test(t, "no end", func(t *T) {
		t.Equal(1, 1)
	})
}
