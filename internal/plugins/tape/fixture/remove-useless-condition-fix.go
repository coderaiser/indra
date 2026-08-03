//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Ok(err)

		t.End()
	})
}
