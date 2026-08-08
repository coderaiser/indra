//go:build ignore

package fixture

import tape "github.com/coderaiser/go-tape"

func f() {
	tape.Equal(1, tape.T{})
}
