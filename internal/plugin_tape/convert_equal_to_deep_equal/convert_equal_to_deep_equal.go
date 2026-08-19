package convert_equal_to_deep_equal

import (
	"go/ast"

	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "Equal: use DeepEqual for slices" }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, __array)":    "__a.DeepEqual(__b, __array)",
		"__a.Equal(__array, __b)":    "__a.DeepEqual(__array, __b)",
		"__a.NotEqual(__b, __array)": "__a.NotDeepEqual(__b, __array)",
	}
}

// Match guards so Equal is only upgraded to DeepEqual when the compared value
// is an actual slice/struct literal (composite literal) or a variable bound to
// one — matching putout's isArray/isObject checks.
func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, __array)": sliceOrStruct,
		"__a.Equal(__array, __b)": sliceOrStruct,
	}
}

// sliceOrStruct reports whether __array is a composite literal or an identifier
// whose binding is a composite literal.
func sliceOrStruct(vars Vars, path Path) bool {
	node := vars["__array"]
	if IsCompositeLit(node) {
		return true
	}
	if !IsIdent(node) {
		return false
	}
	binding := GetBinding(path, Extract(node))
	if binding == nil {
		return false
	}
	assign, ok := binding.Node.(*ast.AssignStmt)
	if !ok || len(assign.Rhs) == 0 {
		return false
	}
	return IsCompositeLit(assign.Rhs[0])
}

// Plugin wraps the rule for the registry: a replacer with a Match guard. The
// [match] config already scopes tape rules to *_test.go files, so no per-plugin
// import guard is needed.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
