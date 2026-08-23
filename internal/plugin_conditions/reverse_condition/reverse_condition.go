// Package reverse_condition reverses negated comparisons and applies De
// Morgan's identities: !(a > b) becomes a <= b, !(a || b) becomes !a && !b.
// A pure port of putout's reverse-condition with Go's == / != operators; the
// JS guard against double negation is unnecessary because indra only rewrites
// the outer negation of a binary expression, which never produces one.
package reverse_condition

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "reverse condition" }

// unparen strips enclosing parentheses from an expression.

// reversible reports whether the operator of a negated binary expression has
// a rewrite: relational operators flip to their complement, logical ones
// expand through De Morgan's identities.

// not wraps e in a syntactic negation.

// Fix flips relational operators to their complements and expands logical
// operators via De Morgan. The rewrite is idempotent — the result never
// matches again — so the fix loop converges after one pass.
func Fix(p Path, _ map[string]any) {
	u := p.Node.(*ast.UnaryExpr)
	b := unparen(u.X).(*ast.BinaryExpr)

	var repl ast.Expr
	switch b.Op {
	case token.GTR:
		repl = &ast.BinaryExpr{X: b.X, Op: token.LEQ, Y: b.Y}
	case token.GEQ:
		repl = &ast.BinaryExpr{X: b.X, Op: token.LSS, Y: b.Y}
	case token.LSS:
		repl = &ast.BinaryExpr{X: b.X, Op: token.GEQ, Y: b.Y}
	case token.LEQ:
		repl = &ast.BinaryExpr{X: b.X, Op: token.GTR, Y: b.Y}
	case token.LOR:
		repl = &ast.BinaryExpr{X: not(b.X), Op: token.LAND, Y: not(b.Y)}
	case token.LAND:
		repl = &ast.BinaryExpr{X: not(b.X), Op: token.LOR, Y: not(b.Y)}
	}

	p.Replace(repl)
}

func unparen(e ast.Expr) ast.Expr {
	for {
		paren, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = paren.X
	}
}

func reversible(op token.Token) bool {
	switch op {
	case token.GTR, token.GEQ, token.LSS, token.LEQ, token.LOR, token.LAND:
		return true
	}
	return false
}

func findNegatedCondition(p Path, push func(Path)) {
	u := p.Node.(*ast.UnaryExpr)
	if u.Op != token.NOT {
		return
	}
	b, ok := unparen(u.X).(*ast.BinaryExpr)
	if !ok || !reversible(b.Op) {
		return
	}
	push(p)
}

func not(e ast.Expr) ast.Expr {
	return &ast.UnaryExpr{Op: token.NOT, X: e}
}
func Traverse() Traverser {
	return Traverser{"*ast.UnaryExpr": findNegatedCondition}
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
