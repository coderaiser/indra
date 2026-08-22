// Package apply_dedent removes unnecessary dedent.Dedent wrappers around raw
// strings, leaving the string literal in place for the formatter to normalize.
package apply_dedent

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "apply dedent" }

// Fix removes the dedent.Dedent(...) wrapper, keeping the string literal
// verbatim. A later remove-unused-import pass drops the now-unused dedent
// import. p.Node is the pushed CallExpr; nothing here needs options.
func Fix(p Path, _ map[string]any) {
	call := p.Node.(*ast.CallExpr)
	p.Replace(call.Args[0])
}

func Traverse() Traverser {
	return Traverser{"*ast.CallExpr": findDedent}
}

// findDedent pushes a CallExpr when it is a dedent.Dedent call whose single
// argument is a string literal, wherever it appears in the tree.
func findDedent(p Path, push func(Path)) {
	call := p.Node.(*ast.CallExpr)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	recv, ok := sel.X.(*ast.Ident)
	if !ok || recv.Name != "dedent" || sel.Sel.Name != "Dedent" {
		return
	}
	if len(call.Args) != 1 {
		return
	}
	if _, ok := call.Args[0].(*ast.BasicLit); !ok {
		return
	}
	push(p)
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
