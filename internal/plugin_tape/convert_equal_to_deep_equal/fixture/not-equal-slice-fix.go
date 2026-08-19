//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

func f() {
	t.NotDeepEqual(x, []Block{})

}
