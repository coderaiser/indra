package remove_only

import (
	"go/ast"
	"slices"

	"coderaiser/indra/types"
)

func Report() string { return `Remove "Test.Only"` }

func Replace() types.Replacer {
	return types.Replacer{
		`__a.Only(__b, __c, func(__d *__a.T) { __body })`:      "__a(__b, __c, func(__d *__a.T) {\n__body\n})",
		`__a.Only(__b, __c, func(__d *__a.T) { __body }, __e)`: "__a(__b, __c, func(__d *__a.T) {\n__body\n}, __e)",
	}
}

// Filter guards on the receiver name: Test is always allowed, and extra
// allowed receivers come from rule options (allowed = [...] in .indra.toml).
func Filter() types.Filter {
	return types.Filter{
		`__a.Only(__b, __c, func(__d *__a.T) { __body })`:      allowedReceiver,
		`__a.Only(__b, __c, func(__d *__a.T) { __body }, __e)`: allowedReceiver,
	}
}

func allowedReceiver(vars types.Vars, _ types.Path, opts types.Options) bool {
	ident, ok := vars["__a"].(*ast.Ident)
	allowed := append([]string{"Test"}, opts.StringSlice("allowed")...)
	return ok && slices.Contains(allowed, ident.Name)
}

// Plugin wraps the rule for the registry: a replacer whose Filter guard is
// option-aware.
type Plugin struct{}

func (Plugin) Report() string          { return Report() }
func (Plugin) Filter() types.Filter    { return Filter() }
func (Plugin) Replace() types.Replacer { return Replace() }
