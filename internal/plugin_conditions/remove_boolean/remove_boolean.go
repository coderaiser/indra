// Package remove_boolean removes comparisons against boolean literals:
// x == true becomes x, x == false becomes !x, and so on. Go has no Boolean()
// constructor, so unlike putout's version only the comparison subset applies.
package remove_boolean

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove boolean" }

// findBooleanComparison matches == / != comparisons whose right operand is
// the predeclared true or false identifier.

// isBoolLit reports whether n is the predeclared true or false identifier.

// Fix rewrites the comparison per the truth table: comparing against true is
// the operand itself, against false its negation — with NEQ inverting both.
func Fix(p Path, _ map[string]any) {
	e := p.Node.(*ast.BinaryExpr)

	negate := e.Y.(*ast.Ident).Name == "false"
	if e.Op == token.NEQ {
		negate = !negate
	}

	var repl ast.Expr = e.X
	if negate {
		repl = &ast.UnaryExpr{Op: token.NOT, X: e.X}
	}

	p.Replace(repl)
}

func findBooleanComparison(p Path, push func(Path)) {
	e := p.Node.(*ast.BinaryExpr)
	if e.Op != token.EQL && e.Op != token.NEQ {
		return
	}
	if !isBoolLit(e.Y) {
		return
	}
	push(p)
}

func isBoolLit(n ast.Expr) bool {
	id, ok := n.(*ast.Ident)
	return ok && (id.Name == "true" || id.Name == "false")
}
func Traverse() Traverser {
	return Traverser{"*ast.BinaryExpr": findBooleanComparison}
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
