//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

const disabled = false

func f() {
	t.Equal(result, disabled)
}
