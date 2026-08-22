//go:build ignore

package fixture

import . "coderaiser/indra/types"

var unused = 1

func helper() {}

func Traverse() Traverser {
	return Traverser{}
}

func Extra() {}

func Report(_ Path) string { return "reorder" }

func Fix(_ Path, _ map[string]any) {}
