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
	f()
	a.b.Inspect()
	ast.Print(p.Node)
	_ = push
}

// cover the returnsTraverser negative branches: no results, a non-Traverser
// ident result, and a non-ident result type.
func helper() {}

func other() string { return "x" }

func list() []string { return nil }
