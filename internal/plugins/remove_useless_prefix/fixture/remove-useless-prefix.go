//go:build ignore

package fixture

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	tape.Test(t, "foo: bar", func(t *tape.T) {
		t.Equal(1, 1)
		t.End()
	})
}
