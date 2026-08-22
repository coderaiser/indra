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

// Filter guards on the receiver name: Test is always allowed, and extra
// allowed receivers come from rule options (allowed = [...] in .indra.toml).
// The loader binds rule Options into the synthesized matcher, so only Test
// calls are rewritten unless a project opts in to more receivers.
func Filter() types.Filter {
	return types.Filter{
		`__a.Skip(__b, __c, func(__d *__a.T) { __body })`:      allowedReceiver,
		`__a.Skip(__b, __c, func(__d *__a.T) { __body }, __e)`: allowedReceiver,
	}
}

func allowedReceiver(vars types.Vars, _ types.Path, opts types.Options) bool {
	ident, ok := vars["__a"].(*ast.Ident)
	allowed := append([]string{"Test"}, opts.StringSlice("allowed")...)
	return ok && slices.Contains(allowed, ident.Name)
}

// Plugin wraps the rule for the registry: a replacer whose Filter guard is
// option-aware. The [match] config already scopes tape rules to *_test.go
// files, so no per-plugin import guard is needed.
type Plugin struct{}

func (Plugin) Report() string          { return Report() }
func (Plugin) Filter() types.Filter    { return Filter() }
func (Plugin) Replace() types.Replacer { return Replace() }
