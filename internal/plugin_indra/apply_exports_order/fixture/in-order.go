//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report(_ Path) string { return "ordered" }

func Fix(_ Path, _ map[string]any) {}

func Traverse() Traverser {
	return Traverser{}
}
