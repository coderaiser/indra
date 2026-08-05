package extract_result_from_assertion

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

// noResultInBlock is a guard that rejects re-extraction when a "result"
// variable is already declared in the containing block (which would shadow the
// injected declaration). A nil block (no block context) passes.
func noResultInBlock(_ Vars, block *ast.BlockStmt) bool {
	if block == nil {
		return true
	}
	return !blockDeclares(block, "result")
}

// blockDeclares reports whether any statement in block declares name via a
// short variable declaration (:=).
func blockDeclares(block *ast.BlockStmt, name string) bool {
	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			continue
		}
		for _, lhs := range assign.Lhs {
			if id, ok := lhs.(*ast.Ident); ok && id.Name == name {
				return true
			}
		}
	}
	return false
}
