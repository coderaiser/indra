//go:build ignore

package fixture

import (
	"go/ast"

	"coderaiser/indra/compare"
)

type Plugin struct{}

func Report(_ string) string { return "rule" }

func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if compare.GetTemplateValues(s, "__.End()") != nil {
			return true
		}
	}
	return false
}
