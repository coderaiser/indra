//go:build ignore

package fixture

import (
	"go/ast"

	. "coderaiser/indra/types"
)

// mixed-guards has one nil guard and one real guard, so the Matcher is useful.
func Match() Matcher {
	return Matcher{
		`Test(__a, __b, func(__a *T) { __body })`: nil,
		`Test.Only(__a, __b, func(__a *T) { __body })`: func(vars Vars, _ *ast.BlockStmt) bool {
			return vars["__a"] != nil
		},
	}
}
