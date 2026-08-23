// Package merge_if_statements collapses nested ifs into a single if with a
// conjunction: if a { if b { x } } becomes if a && b { x }. The rewrite only
// applies when the outer if has no else and no init statement, its body holds
// exactly one statement, and that statement is an else-less if without an init
// — an init clause would be dropped by the merge. Ported from putout's
// merge-if-statements.
package merge_if_statements

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "merge if statements" }

// Fix folds the inner condition into the outer one and adopts the inner body.
// Positions in the enclosing function are stripped afterwards: the merged
// nodes carry positions from different depths of the file, which makes
// go/printer insert spurious line breaks.
func Fix(p Path, _ map[string]any) {
	outer := p.Node.(*ast.IfStmt)
	inner := outer.Body.List[0].(*ast.IfStmt)
	outer.Cond = &ast.BinaryExpr{
		X:  outer.Cond,
		Op: token.LAND,
		Y:  inner.Cond,
	}
	outer.Body = inner.Body
	if block := enclosingFuncBody(p); block != nil {
		stripPositions(block)
	}
}
func findMergeCandidate(p Path, push func(Path)) {
	outer := p.Node.(*ast.IfStmt)
	if outer.Else != nil || outer.Init != nil {
		return
	}
	if len(outer.Body.List) != 1 {
		return
	}
	inner, ok := outer.Body.List[0].(*ast.IfStmt)
	if !ok || inner.Else != nil || inner.Init != nil {
		return
	}
	push(p)
}
func Traverse() Traverser {
	return Traverser{"*ast.IfStmt": findMergeCandidate}
}

// enclosingFuncBody returns the body block of the outermost enclosing
// function declaration, or nil when there is none.
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
