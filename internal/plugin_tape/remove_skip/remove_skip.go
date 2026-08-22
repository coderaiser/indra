package remove_skip

import (
	"go/ast"
	"slices"

	"coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "remove Test.Skip call" }

func Replace() types.Replacer {
	return types.Replacer{
		`__a.Skip(__b, __c, func(__d *__a.T) { __body })`:      "__a(__b, __c, func(__d *__a.T) {\n__body\n})",
		`__a.Skip(__b, __c, func(__d *__a.T) { __body }, __e)`: "__a(__b, __c, func(__d *__a.T) {\n__body\n}, __e)",
	}
}

// Match builds the guard map once, closing each guard over the rule's Options
// (putout's match({options})). Test is always an allowed receiver; extra
// allowed receivers come from rule options (allowed = [...] in .indra.toml).
// The loader calls Match with the rule's Options, so only Test calls are
// rewritten unless a project opts in to more receivers.
func Match(opts types.Options) types.Matcher {
	return types.Matcher{
		`__a.Skip(__b, __c, func(__d *__a.T) { __body })`:      allowedReceiver(opts),
		`__a.Skip(__b, __c, func(__d *__a.T) { __body }, __e)`: allowedReceiver(opts),
	}
}

// allowedReceiver returns a guard that accepts a Skip when the receiver is an
// allowed name. The base set is always "Test"; opts.allowed extends it.
func allowedReceiver(opts types.Options) types.MatchFn {
	return func(vars types.Vars, _ types.Path) bool {
		ident, ok := vars["__a"].(*ast.Ident)
		allowed := append([]string{"Test"}, opts.StringSlice("allowed")...)
		return ok && slices.Contains(allowed, ident.Name)
	}
}

// Plugin wraps the rule for the registry: a replacer whose Match guard is
// option-aware. The [match] config already scopes tape rules to *_test.go
// files, so no per-plugin import guard is needed.
type Plugin struct{}

func (Plugin) Report() string               { return Report() }
func (Plugin) Match(opts types.Options) types.Matcher { return Match(opts) }
func (Plugin) Replace() types.Replacer      { return Replace() }
