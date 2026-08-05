package extract_result_from_assertion

import (
	"go/ast"
	"go/token"
	"testing"

	. "coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

func TestBlockDeclares(t *testing.T) {
	tape.Test(t, "blockDeclares: returns true when name is declared", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("result")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		t.Equal(blockDeclares(block, "result"), true)
		t.End()
	})

	tape.Test(t, "blockDeclares: returns false when name not declared", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("other")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		t.Equal(blockDeclares(block, "result"), false)
		t.End()
	})

	tape.Test(t, "blockDeclares: non-define assign returns false", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent("result")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		t.Equal(blockDeclares(block, "result"), false)
		t.End()
	})

	tape.Test(t, "blockDeclares: non-assign stmt returns false", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.ExprStmt{X: ast.NewIdent("result")},
		}}
		t.Equal(blockDeclares(block, "result"), false)
		t.End()
	})

	tape.Test(t, "blockDeclares: non-ident lhs returns false", func(t *tape.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{&ast.IndexExpr{X: ast.NewIdent("m"), Index: ast.NewIdent("k")}},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		t.Equal(blockDeclares(block, "result"), false)
		t.End()
	})
}

func TestMatchGuardNilBlock(t *testing.T) {
	tape.Test(t, "Match guard: missing $block key returns true", func(t *tape.T) {
		var guard MatchFn
		for _, g := range Match() {
			guard = g
			break
		}
		t.Equal(guard(Vars{}), true)
		t.End()
	})

	tape.Test(t, "Match guard: non-BlockStmt $block returns true", func(t *tape.T) {
		var guard MatchFn
		for _, g := range Match() {
			guard = g
			break
		}
		t.Equal(guard(Vars{"$block": ast.NewIdent("x")}), true)
		t.End()
	})
}
