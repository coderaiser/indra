package test

import (
	"go/ast"
	"strings"
	"testing"

	loader "coderaiser/indra/engine_loader"
	remove_skip "coderaiser/indra/internal/plugin_tape/remove_skip"
	indratest "coderaiser/indra/test"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

// ── indraLint ────────────────────────────────────────────────────────────────

func TestIndraLintSuccess(t *testing.T) {
	src := []byte("package p\n\nfunc f() {\n\tTest.Skip(t, \"foo: something\", func(t *Test.T) {\n\t\tt.End()\n\t})\n}\n")
	plugins := []any{indratest.PluginArg{Rule: "remove-skip", Plugin: remove_skip.Plugin{}}}

	Test(t, "indraLint: fixes without error", func(tt *T) {
		out, err := indraLint(src, true, plugins)
		tt.Ok(err == nil && strings.Contains(string(out.Out), "Test("))
		tt.End()
	})

	Test(t, "indraLint: output drops Skip", func(tt *T) {
		out, _ := indraLint(src, true, plugins)
		tt.NotOk(strings.Contains(string(out.Out), "Skip"))

		tt.End()
	})
}

func TestIndraLintParseError(t *testing.T) {
	Test(t, "indraLint: returns parse error", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "remove-skip", Plugin: remove_skip.Plugin{}}}
		_, err := indraLint([]byte("package p\nfunc (\n"), false, plugins)
		tt.Ok(err)

		tt.End()
	})
}

// ── validatePlugin ───────────────────────────────────────────────────────────

type synthReplacer struct {
	report  string
	match   types.Matcher
	replace types.Replacer
}

func (s synthReplacer) Report() string          { return s.report }
func (s synthReplacer) Match() types.Matcher    { return s.match }
func (s synthReplacer) Replace() types.Replacer { return s.replace }

type synthTraverser struct{}

func (synthTraverser) Report(_ types.Path) string { return "t" }
func (synthTraverser) Traverse() types.Traverser {
	return types.Traverser{"*ast.File": func(types.Path, func(types.Path)) {}}
}
func (synthTraverser) Fix(_ types.Path, _ map[string]any) {}

func catchPanic(fn func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = strings.TrimSpace(r.(string))
		}
	}()
	fn()
	return ""
}

func TestValidatePluginNilGuard(t *testing.T) {
	Test(t, "validatePlugin: panics on nil MatchFn", func(tt *T) {
		pf := loader.PluginFuncs{Name: "nil-guard", Plugin: synthReplacer{report: "x", match: types.Matcher{"p": nil}, replace: types.Replacer{"p": "q"}}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.Ok(strings.Contains(msg, "nil MatchFn"))
		tt.End()
	})
}

func TestValidatePluginOrphanKey(t *testing.T) {
	Test(t, "validatePlugin: panics on orphan Match key", func(tt *T) {
		pf := loader.PluginFuncs{Name: "orphan-key", Plugin: synthReplacer{report: "x", match: types.Matcher{"p": func(types.Vars, *ast.BlockStmt) bool { return true }}, replace: types.Replacer{}}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.Ok(strings.Contains(msg, "Match key not in Replace"))
		tt.End()
	})
}

func TestValidatePluginTraverserNoPanic(t *testing.T) {
	Test(t, "validatePlugin: no panic for traverser plugin", func(tt *T) {
		pf := loader.PluginFuncs{Name: "trav", Plugin: synthTraverser{}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.Equal(msg, "")
		tt.End()
	})
}
