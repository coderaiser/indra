//go:build ignore

package fixture

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: bar", func(t *T) {
		t.Equal(1, 1)
		t.End()
	})
}
