//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

// add-t-end is the canonical happy path: a Test body missing t.End().
func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
	})
}
