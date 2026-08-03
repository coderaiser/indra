//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test.Skip(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})
}
