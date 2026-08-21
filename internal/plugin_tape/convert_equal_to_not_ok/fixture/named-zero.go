//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

const zero = 0

func f() {
	t.Equal(result, zero)
}
