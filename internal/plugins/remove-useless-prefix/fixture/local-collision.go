//go:build ignore

package fixture

import (
	tape "github.com/coderaiser/go-tape"
)

// T is a locally declared type whose name would collide with tape.T if the
// tape prefix were removed. The rule must skip such files.
type T struct {
	inner *tape.T
}

func New(tt *tape.T) *T {
	return &T{inner: tt}
}
