//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Suite.Skip(t, "allowed-receiver: something", func(t *Suite.T) {
		t.Equal(1, 1)
	})
}
