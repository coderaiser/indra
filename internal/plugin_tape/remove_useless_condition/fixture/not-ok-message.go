//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

func f() {
	t.NotOk(err == nil, "err should be set")
}
