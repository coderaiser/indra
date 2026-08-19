//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

func helper() Block {}

func f() {
	t.Equal(x, helper)
}
