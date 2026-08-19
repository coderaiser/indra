//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "contains-end: something", func(t *Test.T) {
		t.End()
		t.Equal(1, 1)
	})
}
