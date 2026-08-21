package convert_equal_to_not_ok

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

func Report() string { return "convert Equal(x, nil/false) to NotOk(x)" }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, nil)":            "__a.NotOk(__b)",
		"__a.Equal(__b, nil, __c)":       "__a.NotOk(__b, __c)",
		"__a.Equal(__b, false)":          "__a.NotOk(__b)",
		"__a.Equal(__b, false, __c)":     "__a.NotOk(__b, __c)",
		`__a.Equal(__b, "")`:             "__a.NotOk(__b)",
		"__a.DeepEqual(__b, nil)":        "__a.NotOk(__b)",
		"__a.DeepEqual(__b, nil, __c)":   "__a.NotOk(__b, __c)",
		"__a.DeepEqual(__b, false)":      "__a.NotOk(__b)",
		"__a.DeepEqual(__b, false, __c)": "__a.NotOk(__b, __c)",
		"__a.Equal(__b, __c)":            "__a.NotOk(__b)",
		"__a.DeepEqual(__b, __c)":        "__a.NotOk(__b)",
		"__a.Equal(__b)":                 "__a.NotOk(__b)",
		"__a.NotEqual(__b)":              "__a.Ok(__b)",
	}
}

// Match guards the generic three-argument Equal/DeepEqual form against falsy
// *literal* constants, mirroring putout's convert-equal-to-not-ok. Dedicated
// Replace keys above handle the nil/false/empty-string identifier and literal
// forms directly; this guard covers numeric/string literals that reduce to a
// falsy value (e.g. t.Equal(x, 0)).
func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, __c)":     isFalsyLiteralGuard,
		"__a.DeepEqual(__b, __c)": isFalsyLiteralGuard,
	}
}

// isFalsyLiteralGuard matches when __c is a literal that Compute reduces to
// false, the integer 0, or the empty string. It rejects numerics and non-empty
// strings so those fall through to putout's explicit Replace keys or stay
// untouched.
func isFalsyLiteralGuard(vars Vars, p Path) bool {
	return isFalsyLiteral(vars["__c"], p)
}

// isFalsyLiteral reports whether c is a falsy scalar value: a BasicLit that
// Compute reduces to false, 0, or "", or a named constant/identifier that
// resolves to one of those values.
func isFalsyLiteral(c ast.Node, p Path) bool {
	lit, ok := c.(*ast.BasicLit)
	if !ok {
		return isFalsyNamedConst(p, c)
	}
	// Float literals are never treated as falsy here: an integral 0 (e.g.
	// t.Equal(x, 0)) reduces to the int value 0 and is convertible, while a
	// float 0.0 falls through to the generic "no report" path below.
	if lit.Kind == token.FLOAT {
		return false
	}
	ok2, val := Compute(c)
	if !ok2 {
		return false
	}
	if s, ok := val.(string); ok && s != "" {
		return false
	}
	return val == int64(0) || val == false || val == ""
}

// isFalsyNamedConst resolves c (expected to be an *ast.Ident) to its declaration
// and reports whether the declared value is falsy. Built-in identifiers false,
// nil, and true are handled directly; other names are resolved via GetBinding.
func isFalsyNamedConst(p Path, c ast.Node) bool {
	ident, ok := c.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name == "true" {
		return false
	}
	if ident.Name == "false" || ident.Name == "nil" {
		return true
	}
	binding := GetBinding(p, ident.Name)
	if binding == nil {
		return false
	}
	genDecl, ok := binding.Node.(*ast.GenDecl)
	if !ok {
		return false
	}
	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, name := range vs.Names {
			if name.Name == ident.Name && i < len(vs.Values) {
				return isFalsyLiteral(vs.Values[i], p)
			}
		}
	}
	return false
}

type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
