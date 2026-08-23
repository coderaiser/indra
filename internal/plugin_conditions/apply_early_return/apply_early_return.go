// Package apply_early_return converts an if/else that ends a function body
// into an early return: a bare return is appended to the consequent and the
// else statements are hoisted after the if. It is the inverse of
// remove-useless-else — that rule fires when the consequent already ends
// with a return, this one adds it. Ported from putout's apply-early-return.
package apply_early_return

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "apply early return" }

func Traverse() Traverser {
	return Traverser{"*ast.IfStmt": findEarlyReturn}
}

func findEarlyReturn(p Path, push func(Path)) {
	ifStmt := p.Node.(*ast.IfStmt)
	if ifStmt.Else == nil {
		return
	}
	// Only plain else blocks are unwrapped; an else-if chain would change
	// meaning once the added return short-circuits it.
	if _, ok := ifStmt.Else.(*ast.BlockStmt); !ok {
		return
	}
	body := ifStmt.Body.List
	if len(body) > 0 && isReturnLike(body[len(body)-1]) {
		return
	}
	// The if must end its immediate function body — otherwise hoisting a bare return
	// into the consequent would skip the statements that follow it.
	if lastStmt(nearestFuncBody(p)) != ifStmt {
		return
	}
	push(p)
}

// lastStmt returns the final statement of a block, or nil when the block is
// missing or empty.
func lastStmt(block *ast.BlockStmt) ast.Stmt {
	var last ast.Stmt
	if block != nil && len(block.List) > 0 {
		last = block.List[len(block.List)-1]
	}
	return last
}

// nearestFuncBody returns the body of the function the if sits in directly —
// the nearest declaration or literal on the stack.
func nearestFuncBody(p Path) *ast.BlockStmt {
	var body *ast.BlockStmt
	for i := len(p.Stack) - 1; i >= 0 && body == nil; i-- {
		switch fn := p.Stack[i].(type) {
		case *ast.FuncDecl:
			body = fn.Body
		case *ast.FuncLit:
			body = fn.Body
		}
	}
	return body
}

// enclosingFuncBody returns the body of the outermost enclosing function —
// the declaration when there is one, the outermost literal otherwise. The
// whole subtree is normalized so no stale positions remain on any sibling.
func enclosingFuncBody(p Path) *ast.BlockStmt {
	var body *ast.BlockStmt
	for i := 0; i < len(p.Stack); i++ {
		switch fn := p.Stack[i].(type) {
		case *ast.FuncDecl:
			body = fn.Body
		case *ast.FuncLit:
			if body == nil {
				body = fn.Body
			}
		}
	}
	return body
}

// isReturnLike reports whether s terminates control flow: a return, break or
// continue.
func isReturnLike(s ast.Stmt) bool {
	switch s := s.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return s.Tok == token.BREAK || s.Tok == token.CONTINUE
	}
	return false
}

// Fix appends a bare return to the consequent and hoists every else statement
// after the if (inserting in reverse so each lands directly below the
// previous one), then drops the else clause. Positions in the enclosing
// function are stripped afterwards to keep go/printer from inserting
// spurious blank lines around the moved statements.
func Fix(p Path, _ map[string]any) {
	ifStmt := p.Node.(*ast.IfStmt)
	ifStmt.Body.List = append(ifStmt.Body.List, &ast.ReturnStmt{})
	elseBlock := ifStmt.Else.(*ast.BlockStmt)
	for i := len(elseBlock.List) - 1; i >= 0; i-- {
		p.InsertAfter(elseBlock.List[i])
	}
	ifStmt.Else = nil
	stripPositions(enclosingFuncBody(p))
}

// stripPositions zeroes every token.Pos field in the sub-tree rooted at root.
func stripPositions(root ast.Node) {
	posType := reflect.TypeOf(token.Pos(0))
	astutil.Apply(root, func(c *astutil.Cursor) bool {
		if c.Node() == nil {
			return true
		}
		v := reflect.ValueOf(c.Node()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Type() == posType {
				v.Field(i).SetInt(0)
			}
		}
		return true
	}, nil)
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
