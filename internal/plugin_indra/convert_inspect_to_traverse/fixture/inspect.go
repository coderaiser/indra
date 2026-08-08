//go:build ignore

package fixture

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "rule" }

func Traverse() Traverser {
	return Traverser{"*ast.File": find}
}

func find(p Path, push func(Path)) {
	ast.Inspect(p.Node, func(n ast.Node) bool { return true })
	_ = push
}
