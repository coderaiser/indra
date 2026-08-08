//go:build ignore

package fixture

import tape "github.com/coderaiser/go-tape"

// Equal is a locally declared variable whose name would collide with tape.Equal
// if the tape prefix were removed. The rule must skip such files.
var Equal int

func f() {
	_ = tape.Equal
}
