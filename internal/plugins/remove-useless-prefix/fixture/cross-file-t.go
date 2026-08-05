//go:build ignore

package fixture

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

// T is used as a bare identifier elsewhere in this file, so removing the tape
// prefix would create an ambiguity between tape.T and the bare T usage.
func TestFoo(t *testing.T) {
	tape.Test(t, "foo", func(t *tape.T) {
		t.End()
	})
}

func usesBareT(_ T) {}
