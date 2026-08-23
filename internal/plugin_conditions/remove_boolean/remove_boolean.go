// Package remove_boolean removes comparisons against boolean literals:
// x == true becomes x, x == false becomes !x, and so on. Go has no Boolean()
// constructor, so unlike putout's version only the comparison subset applies.
//
// indra works on the bare AST, so the rewrite is guarded: the compared operand
// must be statically boolean — either a boolean-shaped expression (!a, a && b,
// a == b) or an identifier whose visible declaration says bool. Rewriting an
// any-typed comparison like val == false to !val would not compile.
package remove_boolean

import (
	"go/ast"
	"go/token"

	"coderaiser/indra/operator"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove boolean" }

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
		repl = negateExpr(e.X)
	}
	if paren, ok := repl.(*ast.ParenExpr); ok {
		repl = paren.X
	}

	p.Replace(repl)
}

func Traverse() Traverser {
	return Traverser{"*ast.BinaryExpr": findBooleanComparison}
}

// findBooleanComparison matches == / != comparisons against true or false
// whose left operand is statically boolean.
func findBooleanComparison(p Path, push func(Path)) {
	e := p.Node.(*ast.BinaryExpr)
	if e.Op != token.EQL && e.Op != token.NEQ {
		return
	}
	if !isBoolLit(e.Y) {
		return
	}
	if !isBoolOperand(p, e.X) {
		return
	}
	push(p)
}

// isBoolLit reports whether n is the predeclared true or false identifier.
func isBoolLit(n ast.Expr) bool {
	id, ok := n.(*ast.Ident)
	return ok && (id.Name == "true" || id.Name == "false")
}

// isBoolOperand reports whether the compared expression is statically boolean:
// a boolean-shaped expression, or an identifier declared bool.
func isBoolOperand(p Path, e ast.Expr) bool {
	if exprIsBool(e) {
		return true
	}
	id, ok := e.(*ast.Ident)
	if !ok {
		return false
	}
	return identIsBool(p, id)
}

// exprIsBool reports whether e's shape guarantees a bool result.
func exprIsBool(e ast.Expr) bool {
	found := false
	switch e := e.(type) {
	case *ast.UnaryExpr:
		found = e.Op == token.NOT
	case *ast.BinaryExpr:
		switch e.Op {
		case token.LAND, token.LOR, token.EQL, token.NEQ,
			token.LSS, token.GTR, token.LEQ, token.GEQ:
			found = true
		}
	case *ast.ParenExpr:
		found = exprIsBool(e.X)
	}
	return found
}

// identIsBool reports whether the visible declaration of id states bool.
func identIsBool(p Path, id *ast.Ident) bool {
	bind := operator.GetBinding(p, id.Name)
	if bind == nil {
		return false
	}
	declared := false
	switch n := bind.Node.(type) {
	case *ast.FuncDecl:
		declared = fieldTypeIs(n.Type.Params, id.Name) ||
			fieldTypeIs(n.Type.Results, id.Name)
	case *ast.FuncLit:
		declared = fieldTypeIs(n.Type.Params, id.Name)
	case *ast.AssignStmt:
		declared = defineWithBoolLiteral(n, id.Name)
	case *ast.DeclStmt:
		if gd, ok := n.Decl.(*ast.GenDecl); ok {
			declared = specsSayBool(gd.Specs, id.Name)
		}
	case *ast.GenDecl:
		declared = specsSayBool(n.Specs, id.Name)
	}
	return declared
}

// fieldTypeIs reports whether the named parameter or result is declared bool.
func fieldTypeIs(fields *ast.FieldList, name string) bool {
	found := false
	for _, f := range fields.List {
		if t, ok := f.Type.(*ast.Ident); ok && t.Name == "bool" {
			for _, n := range f.Names {
				if n.Name == name {
					found = true
				}
			}
		}
	}
	return found
}

// defineWithBoolLiteral reports whether a := statement defines name from the
// predeclared true or false — the only short-declaration form that proves the
// variable's type.
func defineWithBoolLiteral(stmt *ast.AssignStmt, name string) bool {
	defined := false
	for i, lhs := range stmt.Lhs {
		id, ok := lhs.(*ast.Ident)
		if !ok || id.Name != name || i >= len(stmt.Rhs) {
			continue
		}
		if lit, ok := stmt.Rhs[i].(*ast.Ident); ok {
			defined = lit.Name == "true" || lit.Name == "false"
		}
	}
	return defined
}

// specsSayBool reports whether the value spec naming name declares bool —
// either via an explicit bool type or a true/false initializer.
func specsSayBool(specs []ast.Spec, name string) bool {
	declared := false
	for _, spec := range specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, n := range vs.Names {
			if n.Name != name {
				continue
			}
			if t, ok := vs.Type.(*ast.Ident); ok && t.Name == "bool" {
				declared = true
			}
			for _, v := range vs.Values {
				if lit, ok := v.(*ast.Ident); ok {
					declared = declared || lit.Name == "true" || lit.Name == "false"
				}
			}
		}
	}
	return declared
}

// negateExpr negates e, collapsing a double negation (!!e becomes e).
func negateExpr(e ast.Expr) ast.Expr {
	if u, ok := e.(*ast.UnaryExpr); ok && u.Op == token.NOT {
		return u.X
	}
	return &ast.UnaryExpr{Op: token.NOT, X: e}
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
