//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "assign-end: something", func(t *Test.T) {
		result := t.End()
		t.Equal(1, 1)
	})
}
