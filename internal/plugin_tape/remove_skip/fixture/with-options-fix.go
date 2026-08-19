//go:build ignore

package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "with-options: something", func(t *Test.T) {
		t.Equal(1, 1)
	}, testOpts)

}
