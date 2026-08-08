//go:build ignore

package fixture

import (
	"go/ast"

	"coderaiser/indra/compare"
)

func walk(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if compare.GetTemplateValues(s, "__.End()") != nil {
			return true
		}
	}
	return false
}
