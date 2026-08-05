package extract_result_from_assertion

import (
	"go/ast"
	"go/token"
	"testing"

	. "coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

func TestBlockDeclares(t *testing.T) {
	Test(t, "blockDeclares: returns true when name is declared", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("result")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		result := blockDeclares(block, "result")
		t.Equal(result, true)

		t.End()
	})

	Test(t, "blockDeclares: returns false when name not declared", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{ast.NewIdent("other")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		result := blockDeclares(block, "result")
		t.Equal(result, false)

		t.End()
	})

	Test(t, "blockDeclares: non-define assign returns false", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent("result")},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		result := blockDeclares(block, "result")
		t.Equal(result, false)

		t.End()
	})

	Test(t, "blockDeclares: non-assign stmt returns false", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.ExprStmt{X: ast.NewIdent("result")},
		}}
		result := blockDeclares(block, "result")
		t.Equal(result, false)

		t.End()
	})

	Test(t, "blockDeclares: non-ident lhs returns false", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{&ast.IndexExpr{X: ast.NewIdent("m"), Index: ast.NewIdent("k")}},
				Rhs: []ast.Expr{ast.NewIdent("x")},
			},
		}}
		result := blockDeclares(block, "result")
		t.Equal(result, false)

		t.End()
	})
}

func TestMatchGuardNilBlock(t *testing.T) {
	Test(t, "Match guard: missing $block key returns true", func(t *T) {
		var guard MatchFn
		for _, g := range Match() {
			guard = g
			break
		}
		result := guard(Vars{})
		t.Equal(result, true)

		t.End()
	})

	Test(t, "Match guard: non-BlockStmt $block returns true", func(t *T) {
		var guard MatchFn
		for _, g := range Match() {
			guard = g
			break
		}
		result := guard(Vars{"$block": ast.NewIdent("x")})
		t.Equal(result, true)

		t.End()
	})
}
