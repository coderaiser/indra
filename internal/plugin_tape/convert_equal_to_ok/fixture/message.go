//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

func f() {
	t.Equal(result, true, "should be true")
}