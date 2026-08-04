//go:build ignore

package fixture

import (
	. "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *T) {
		t.NotOk(err)

		t.End()
	})
}
