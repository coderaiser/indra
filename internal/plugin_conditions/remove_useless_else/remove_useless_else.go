// Package remove_useless_else unwraps else branches of if statements whose
// consequent already ends with a return, break or continue — control never
// reaches the else in that case, so hoisting its statements after the if is
// equivalent and one nesting level shallower. Ported from putout's
// remove-useless-else.
package remove_useless_else

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove useless else" }

// isReturnLike reports whether s terminates control flow: a return, break or
// continue. fallthrough is deliberately excluded — it must stay attached to
// the switch clause it belongs to.

// Fix moves every else statement after the if (inserting in reverse so each
// InsertAfter lands directly below the previous one, preserving order), then
// drops the else clause.
func Fix(p Path, _ map[string]any) {
	ifStmt := p.Node.(*ast.IfStmt)
	elseBlock := ifStmt.Else.(*ast.BlockStmt)
	for i := len(elseBlock.List) - 1; i >= 0; i-- {
		p.InsertAfter(elseBlock.List[i])
	}
	ifStmt.Else = nil

	// Normalize the enclosing block's source positions: the hoisted statements
	// still carry positions from inside the removed else block, which makes
	// go/printer emit a spurious blank line before the closing brace.
	if block := enclosingFuncBody(p); block != nil {
		stripPositions(block)
	}
}
func findUselessElse(p Path, push func(Path)) {
	ifStmt := p.Node.(*ast.IfStmt)
	if ifStmt.Else == nil {
		return
	}
	if _, ok := ifStmt.Else.(*ast.BlockStmt); !ok {
		return
	}
	body := ifStmt.Body.List
	if len(body) == 0 {
		return
	}
	if !isReturnLike(body[len(body)-1]) {
		return
	}
	push(p)
}

func isReturnLike(s ast.Stmt) bool {
	switch s := s.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return s.Tok == token.BREAK || s.Tok == token.CONTINUE
	}
	return false
}
func Traverse() Traverser {
	return Traverser{"*ast.IfStmt": findUselessElse}
}

// enclosingFuncBody returns the body block of the outermost enclosing function
// declaration, or nil when there is none. Normalizing the whole
// function keeps every position in it mutually consistent — normalizing only
// the direct parent block would leave stale positions on its siblings.
func enclosingFuncBody(p Path) *ast.BlockStmt {
	var body *ast.BlockStmt
	for i := 0; i < len(p.Stack); i++ {
		if fn, ok := p.Stack[i].(*ast.FuncDecl); ok {
			body = fn.Body
			break
		}
	}
	return body
}

// stripPositions zeroes every token.Pos field in the sub-tree rooted at root.
// Generated nodes that reuse original AST nodes would otherwise carry stale
// source positions that make go/printer insert spurious line breaks.
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
