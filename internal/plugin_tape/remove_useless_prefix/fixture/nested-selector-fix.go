//go:build ignore

package fixture

import . "github.com/coderaiser/go-tape"

func f() {
	Equal(1, T{})
}
