package remove_only

import (
	"go/ast"
	"slices"

	"coderaiser/indra/types"
)

func Report() string { return `Remove "Test.Only"` }

// Match builds the guard map once, closing each guard over the rule's Options
// (putout's match({options})). Test is always an allowed receiver; extra
// allowed receivers come from rule options (allowed = [...] in .indra.toml).
func Match(opts types.Options) types.Matcher {
	return types.Matcher{
		`__a.Only(__b, __c, func(__d *__a.T) { __body })`:      allowedReceiver(opts),
		`__a.Only(__b, __c, func(__d *__a.T) { __body }, __e)`: allowedReceiver(opts),
	}
}
func Replace() types.Replacer {
	return types.Replacer{
		`__a.Only(__b, __c, func(__d *__a.T) { __body })`:      "__a(__b, __c, func(__d *__a.T) {\n__body\n})",
		`__a.Only(__b, __c, func(__d *__a.T) { __body }, __e)`: "__a(__b, __c, func(__d *__a.T) {\n__body\n}, __e)",
	}
}

// allowedReceiver returns a guard that accepts an Only when the receiver is an
// allowed name. The base set is always "Test"; opts.allowed extends it.
func allowedReceiver(opts types.Options) types.MatchFn {
	return func(vars types.Vars, _ types.Path) bool {
		ident, ok := vars["__a"].(*ast.Ident)
		allowed := append([]string{"Test"}, opts.StringSlice("allowed")...)
		return ok && slices.Contains(allowed, ident.Name)
	}
}

// Plugin wraps the rule for the registry: a replacer whose Match guard is
// option-aware.
type Plugin struct{}

func (Plugin) Report() string                         { return Report() }
func (Plugin) Match(opts types.Options) types.Matcher { return Match(opts) }
func (Plugin) Replace() types.Replacer                { return Replace() }
