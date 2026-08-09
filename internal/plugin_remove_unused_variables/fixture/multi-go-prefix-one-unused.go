//go:build ignore

package fixture

import (
	"github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-coverage"
)

func F() { tape.Test(nil, "", nil) }
