//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

func f() {
	t.NotOk(err, "err should be set")

}
