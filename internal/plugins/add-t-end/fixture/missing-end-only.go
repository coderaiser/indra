//go:build ignore

package fixture

import (
	"testing"
	Test "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	Test.Only(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
	})
}
