// Package merge_if_with_else collapses an if/else-if pair whose branches are
// identical into a single if with a disjunction: if a { x } else if b { x }
// becomes if a || b { x }. Ported from putout's merge-if-with-else.
package merge_if_with_else

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "merge if with else" }

// Init statements would be dropped by the merge.

// An else-if that carries its own else would lose that branch on merge.

// equalBodies reports whether two blocks print identically. Printing rather
// than comparing nodes sidesteps position differences between the branches.

// printed renders a node the way go/printer does for parsed source.

// Fix folds the else-if condition into the outer one with || and drops the
// else branch. Positions in the enclosing function are stripped afterwards:
// the merged condition mixes nodes from different depths of the file, which
// makes go/printer insert spurious line breaks.
func Fix(p Path, _ map[string]any) {
	ifStmt := p.Node.(*ast.IfStmt)
	elseIf := ifStmt.Else.(*ast.IfStmt)
	ifStmt.Cond = &ast.BinaryExpr{
		X:  ifStmt.Cond,
		Op: token.LOR,
		Y:  elseIf.Cond,
	}
	ifStmt.Else = nil
	stripPositions(enclosingFuncBody(p))
}
func findMergeWithElse(p Path, push func(Path)) {
	ifStmt := p.Node.(*ast.IfStmt)
	elseIf, ok := ifStmt.Else.(*ast.IfStmt)
	if !ok {
		return
	}

	if ifStmt.Init != nil || elseIf.Init != nil {
		return
	}

	if elseIf.Else != nil {
		return
	}
	if !equalBodies(ifStmt.Body, elseIf.Body) {
		return
	}
	push(p)
}

func equalBodies(a, b *ast.BlockStmt) bool {
	return printed(a) == printed(b)
}

func printed(n ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), n)
	return buf.String()
}
func Traverse() Traverser {
	return Traverser{"*ast.IfStmt": findMergeWithElse}
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
