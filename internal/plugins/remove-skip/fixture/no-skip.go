//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})
}
