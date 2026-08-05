//go:build ignore

package fixture

import tape "github.com/coderaiser/go-tape"

func f() {
	x := 1
	foo()
	t.NotOk(other)
	t.NoError(err)
}
