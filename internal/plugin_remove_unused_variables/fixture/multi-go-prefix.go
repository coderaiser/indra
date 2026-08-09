//go:build ignore

package fixture

import (
	"github.com/coderaiser/go-coverage"
	"github.com/coderaiser/go-tape"
)

func F() {
	tape.Test(nil, "", nil)
	coverage.Process("")
}
