//go:build ignore

package fixture

import "go/ast"

func walk(n ast.Node) {
	ast.Inspect(n, func(x ast.Node) bool { return true })
}
