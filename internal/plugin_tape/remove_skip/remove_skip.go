package remove_skip

import (
	"go/ast"

	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "remove Test.Skip call" }

func Replace() Replacer {
	return Replacer{
		`__a.Skip(__b, __c, func(__d *__a.T) { __body })`:      "__a(__b, __c, func(__d *__a.T) {\n__body\n})",
		`__a.Skip(__b, __c, func(__d *__a.T) { __body }, __e)`: "__a(__b, __c, func(__d *__a.T) {\n__body\n}, __e)",
	}
}

// Match guards on the receiver name Test, so only Test.Skip calls (not an
// unrelated Skip method) are rewritten.
func Match() Matcher {
	return Matcher{
		`__a.Skip(__b, __c, func(__d *__a.T) { __body })`:      receiverIsTest,
		`__a.Skip(__b, __c, func(__d *__a.T) { __body }, __e)`: receiverIsTest,
	}
}

func receiverIsTest(vars Vars, _ Path) bool {
	ident, ok := vars["__a"].(*ast.Ident)
	return ok && ident.Name == "Test"
}

// Plugin wraps the rule for the registry: a replacer with a Match guard. The
// [match] config already scopes tape rules to *_test.go files, so no per-plugin
// import guard is needed.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
