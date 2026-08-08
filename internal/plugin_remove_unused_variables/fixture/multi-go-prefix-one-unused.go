//go:build ignore

package fixture

import (
	"github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-coverage"
)

func f() { tape.Test(nil, "", nil) }
