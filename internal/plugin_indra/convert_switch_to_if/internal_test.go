package convert_switch_to_if

import (
	"go/ast"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
	. "coderaiser/indra/types"
)

func switchPath(sw *ast.SwitchStmt) Path {
	return Path{Node: sw}
}

func TestHasReturn(t *testing.T) {
	Test(t, "hasReturn: true when ReturnStmt present", func(t *T) {
		body := &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}}
		p := Path{Node: body}
		t.Ok(hasReturn(p))
		t.End()
	})

	Test(t, "hasReturn: false when no ReturnStmt", func(t *T) {
		body := &ast.BlockStmt{List: []ast.Stmt{
			&ast.ExprStmt{X: ast.NewIdent("x")},
		}}
		p := Path{Node: body}
		t.Equal(hasReturn(p), false)
		t.End()
	})
}

func TestHasFallthrough(t *testing.T) {
	Test(t, "hasFallthrough: true when fallthrough present", func(t *T) {
		body := &ast.BlockStmt{List: []ast.Stmt{
			&ast.BranchStmt{Tok: token.FALLTHROUGH},
		}}
		p := Path{Node: body}
		t.Ok(hasFallthrough(p))
		t.End()
	})

	Test(t, "hasFallthrough: false when no fallthrough", func(t *T) {
		body := &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}}
		p := Path{Node: body}
		t.Equal(hasFallthrough(p), false)
		t.End()
	})
}

func TestRemoveBreaks(t *testing.T) {
	Test(t, "removeBreaks: removes break statements", func(t *T) {
		stmts := []ast.Stmt{
			&ast.ReturnStmt{},
			&ast.BranchStmt{Tok: token.BREAK},
		}
		result := removeBreaks(stmts)
		t.Equal(len(result), 1)
		t.End()
	})

	Test(t, "removeBreaks: keeps non-break statements", func(t *T) {
		stmts := []ast.Stmt{&ast.ReturnStmt{}}
		result := removeBreaks(stmts)
		t.Equal(len(result), 1)
		t.End()
	})
}

func TestSwitchBlock(t *testing.T) {
	Test(t, "switchBlock: nil when no parent", func(t *T) {
		sw := &ast.SwitchStmt{Body: &ast.BlockStmt{}}
		p := Path{Node: sw}
		t.Ok(switchBlock(p) == nil)
		t.End()
	})

	Test(t, "switchBlock: nil when parent is not a block", func(t *T) {
		sw := &ast.SwitchStmt{Body: &ast.BlockStmt{}}
		p := Path{Node: sw, Stack: []ast.Node{ast.NewIdent("x")}}
		t.Ok(switchBlock(p) == nil)
		t.End()
	})

	Test(t, "switchBlock: returns enclosing block", func(t *T) {
		block := &ast.BlockStmt{}
		sw := &ast.SwitchStmt{Body: &ast.BlockStmt{}}
		p := Path{Node: sw, Stack: []ast.Node{block}}
		t.Ok(switchBlock(p) == block)
		t.End()
	})
}

func TestIsConvertible(t *testing.T) {
	Test(t, "isConvertible: false when tag is nil", func(t *T) {
		sw := &ast.SwitchStmt{Body: &ast.BlockStmt{}}
		t.Equal(isConvertible(switchPath(sw)), false)
		t.End()
	})

	Test(t, "isConvertible: true for zero-case tag switch", func(t *T) {
		sw := &ast.SwitchStmt{
			Tag:  ast.NewIdent("x"),
			Body: &ast.BlockStmt{},
		}
		t.Equal(isConvertible(switchPath(sw)), true)
		t.End()
	})

	Test(t, "isConvertible: false when default clause present", func(t *T) {
		sw := &ast.SwitchStmt{
			Tag: ast.NewIdent("x"),
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.CaseClause{Body: []ast.Stmt{&ast.ReturnStmt{}}},
			}},
		}
		t.Equal(isConvertible(switchPath(sw)), false)
		t.End()
	})

	Test(t, "isConvertible: false when case has no return", func(t *T) {
		sw := &ast.SwitchStmt{
			Tag:  ast.NewIdent("x"),
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.CaseClause{List: []ast.Expr{ast.NewIdent(`"a"`)}, Body: []ast.Stmt{
					&ast.ExprStmt{X: ast.NewIdent("println")},
				}},
			}},
		}
		t.Equal(isConvertible(switchPath(sw)), false)
		t.End()
	})

	Test(t, "isConvertible: false when case has fallthrough", func(t *T) {
		sw := &ast.SwitchStmt{
			Tag:  ast.NewIdent("x"),
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.CaseClause{List: []ast.Expr{ast.NewIdent(`"a"`)}, Body: []ast.Stmt{
					&ast.BranchStmt{Tok: token.FALLTHROUGH},
				}},
			}},
		}
		t.Equal(isConvertible(switchPath(sw)), false)
		t.End()
	})

	Test(t, "isConvertible: true for convertible switch", func(t *T) {
		sw := &ast.SwitchStmt{
			Tag:  ast.NewIdent("x"),
			Body: &ast.BlockStmt{List: []ast.Stmt{
				&ast.CaseClause{List: []ast.Expr{ast.NewIdent(`"a"`)}, Body: []ast.Stmt{
					&ast.ReturnStmt{},
				}},
			}},
		}
		t.Equal(isConvertible(switchPath(sw)), true)
		t.End()
	})
}
