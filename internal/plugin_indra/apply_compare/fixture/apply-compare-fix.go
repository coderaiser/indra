//go:build ignore

package fixture

import (
	. "coderaiser/indra/operator"
	"go/ast"
)

type Plugin struct{}

func Report(_ string) string { return "rule" }

func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if Compare(s, "__.End()") {
			return true
		}
	}
	return false
}
