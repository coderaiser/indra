//go:build ignore

package fixture

import (
	"go/ast"

	"coderaiser/indra/compare"
)

type Plugin struct{}

func Report(_ string) string { return "rule" }

// op mismatch: EQ instead of !=
func a(node ast.Node) {
	if compare.GetTemplateValues(node, "__.End()") == nil {
		return
	}
}

// X is not a call
func b(x any) {
	if x != nil {
		return
	}
}

// Fun is not a selector
func c(node ast.Node) {
	if otherCall(node) != nil {
		return
	}
}

// sel.X is not the compare ident
func d(node ast.Node, object any) {
	if object.compare.GetTemplateValues(node, "__.End()") != nil {
		return
	}
}

// sel.Sel is not GetTemplateValues
func e(node ast.Node) {
	if compare.Other(node) != nil {
		return
	}
}

// Y is not nil
func f(node ast.Node) {
	if compare.GetTemplateValues(node, "__.End()") != 1 {
		return
	}
}
