package apply_assertions_order

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "Apply assertions order" }

// Fix swaps the interleaved statement with the assertion before it inside
// the parent block: [assertion, gap, End] becomes [gap, assertion, End].
// The Traverse guard already validated the window shape.
func Fix(path Path, _ map[string]any) {
	parent, _ := path.ParentPath()
	block := parent.Node.(*ast.BlockStmt)
	idx := -1
	for i, stmt := range block.List {
		if stmt == path.Node {
			idx = i
			break
		}
	}
	block.List[idx-2], block.List[idx-1] = block.List[idx-1], block.List[idx-2]
}

// Traverse pushes every <recv>.End() preceded by [<recv> assertion,
// non-assertion, End] — the interleaved statement belongs before the
// assertions so they stay adjacent ahead of End.
func Traverse() Traverser {
	return Traverser{
		`__a.End()`: func(path Path, push func(Path)) {
			recv := assertionRecv(path.Node)
			if recv == nil {
				return
			}
			prev, ok := path.PrevSibling()
			if !ok || prev.Node == nil {
				return
			}
			if isAssertion(prev.Node, recv) {
				return
			}
			prevPrev, ok := prev.PrevSibling()
			if !ok || prevPrev.Node == nil {
				return
			}
			if !isAssertion(prevPrev.Node, recv) {
				return
			}
			push(path)
		},
	}
}

// assertionRecv returns the receiver identifier of a <recv>.<method>(...)
// expression statement, or nil.
func assertionRecv(node ast.Node) *ast.Ident {
	stmt, ok := node.(*ast.ExprStmt)
	if !ok {
		return nil
	}
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	id, _ := sel.X.(*ast.Ident)
	return id
}

// isAssertion reports whether node is a method call on the given receiver.
func isAssertion(node ast.Node, recv *ast.Ident) bool {
	r := assertionRecv(node)
	return r != nil && r.Name == recv.Name
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
