package test

import (
	"strings"
	"testing"

	loader "coderaiser/indra/engine_loader"
	remove_skip "coderaiser/indra/internal/plugin_tape/remove_skip"
	indratest "coderaiser/indra/test"
	"coderaiser/indra/types"
)

// run is a fixture-less runner for white-box tests that call engine
// functions directly rather than through fixture files on disk.
// It uses tape.Extend bound to indraLint so *T carries Ok, Equal, etc.
var run = CreateTest("internal-test", remove_skip.Plugin{})

func TestIndraLintSuccess(t *testing.T) {
	src := []byte("package p\n\nimport \"github.com/coderaiser/go-tape\"\n\nfunc f() {\n\tTest.Skip(t, \"foo: something\", func(t *Test.T) {\n\t\tt.End()\n\t})\n}\n")
	plugins := []any{indratest.PluginArg{Rule: "remove-skip", Plugin: remove_skip.Plugin{}}}

	run(t, "internal-test: indraLint: fixes without error", func(tt *T) {
		out, err := indraLint(src, true, plugins)
		tt.Ok(err == nil && strings.Contains(string(out.Out), "Test("))
		tt.End()
	})

	run(t, "internal-test: indraLint: output drops Skip", func(tt *T) {
		out, _ := indraLint(src, true, plugins)
		tt.NotOk(strings.Contains(string(out.Out), "Skip"))
		tt.End()
	})
}

func TestIndraLintParseError(t *testing.T) {
	run(t, "internal-test: indraLint: invalid Go with replacer returns no error", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "remove-skip", Plugin: remove_skip.Plugin{}}}
		_, err := indraLint([]byte("package p\nfunc (\n"), false, plugins)
		tt.NotOk(err)
		tt.End()
	})

	run(t, "internal-test: indraLint: text fallback reports no places on no match", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "remove-skip", Plugin: remove_skip.Plugin{}}}
		out, _ := indraLint([]byte("package p\nfunc (\n"), false, plugins)
		tt.NotOk(len(out.Places))

		tt.End()
	})

	run(t, "internal-test: indraLint: returns parse error for non-replacer plugins", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "synth", Plugin: synthTraverser{}}}
		_, err := indraLint([]byte("package p\nfunc (\n"), false, plugins)
		tt.Ok(err)
		tt.End()
	})
}

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
	run(t, "internal-test: validatePlugin: panics on nil MatchFn", func(tt *T) {
		pf := loader.PluginFuncs{Name: "nil-guard", Plugin: synthReplacer{report: "x", match: types.Matcher{"p": nil}, replace: types.Replacer{"p": "q"}}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.Ok(strings.Contains(msg, "nil MatchFn"))
		tt.End()
	})
}

func TestValidatePluginOrphanKey(t *testing.T) {
	run(t, "internal-test: validatePlugin: panics on orphan Match key", func(tt *T) {
		pf := loader.PluginFuncs{Name: "orphan-key", Plugin: synthReplacer{report: "x", match: types.Matcher{"p": func(types.Vars, types.Path) bool { return true }}, replace: types.Replacer{}}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.Ok(strings.Contains(msg, "Match key not in Replace"))
		tt.End()
	})
}

func TestValidatePluginTraverserNoPanic(t *testing.T) {
	run(t, "internal-test: validatePlugin: no panic for traverser plugin", func(tt *T) {
		pf := loader.PluginFuncs{Name: "trav", Plugin: synthTraverser{}}
		msg := catchPanic(func() { validatePlugin(loader.Load([]loader.PluginFuncs{pf}, loader.Config{})[0]) })
		tt.NotOk(msg)

		tt.End()
	})
}

func TestCreateTestConfigOptions(t *testing.T) {
	suite := CreateTestConfig("remove-skip", remove_skip.Plugin{}, loader.Config{
		"remove-skip": {Enabled: true, Options: map[string]any{"allowed": []string{"Suite"}}},
	})

	suite(t, "internal-test: CreateTestConfig: allowed receiver reports", func(tt *T) {
		tt.Report("suite-skip", "remove Test.Skip call")
		tt.End()
	})
}

var textReplacer = synthReplacer{
	report:  "normalise version field",
	replace: types.Replacer{"\"version\": \"__a\"": "\"version\": \"0.0.0\""},
}

func TestTextFallbackReport(t *testing.T) {
	run(t, "internal-test: text fallback: reports a place without fix", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "rename-version", Plugin: textReplacer}}
		out, _ := indraLint([]byte("{\"version\": \"1.2.3\"}\n"), false, plugins)
		tt.Ok(len(out.Places) > 0)
		tt.End()
	})
}

func TestTextFallbackFix(t *testing.T) {
	run(t, "internal-test: text fallback: rewrites matched text on fix", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "rename-version", Plugin: textReplacer}}
		out, _ := indraLint([]byte("{\"version\": \"1.2.3\"}\n"), true, plugins)
		tt.NotOk(strings.Contains(string(out.Out), "1.2.3"))
		tt.End()
	})
}

func TestTextFallbackIdentity(t *testing.T) {
	identity := synthReplacer{
		report:  "identity",
		replace: types.Replacer{"\"version\": \"__a\"": "\"version\": \"__a\""},
	}
	run(t, "internal-test: text fallback: identity rewrite reports nothing", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "identity", Plugin: identity}}
		out, _ := indraLint([]byte("{\"version\": \"1.2.3\"}\n"), false, plugins)
		count := len(out.Places)
		tt.NotOk(count != 0)
		tt.End()
	})
}

func TestTextFallbackSkipsNonReplacer(t *testing.T) {
	run(t, "internal-test: text fallback: skips non-replacer plugins", func(tt *T) {
		plugins := []any{
			indratest.PluginArg{Rule: "synth", Plugin: synthTraverser{}},
			indratest.PluginArg{Rule: "rename-version", Plugin: textReplacer},
		}
		out, _ := indraLint([]byte("{}\n"), false, plugins)
		count := len(out.Places)
		tt.NotOk(count != 0)
		tt.End()
	})
}

func TestTextFallbackExtraPlaceholder(t *testing.T) {
	extra := synthReplacer{
		report:  "extra placeholder",
		replace: types.Replacer{"v": "w__a"},
	}
	run(t, "internal-test: text fallback: template placeholder beyond captures is empty", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "extra", Plugin: extra}}
		out, _ := indraLint([]byte("v x\n"), true, plugins)
		tt.Ok(strings.Contains(string(out.Out), "w"))
		tt.End()
	})
}

func TestTextFallbackPositionSecondLine(t *testing.T) {
	run(t, "internal-test: text fallback: reports position on later lines", func(tt *T) {
		plugins := []any{indratest.PluginArg{Rule: "rename-version", Plugin: textReplacer}}
		out, _ := indraLint([]byte("{\n\"version\": \"1.2.3\"\n}\n"), false, plugins)
		tt.Equal(out.Places[0].Position.Line, 2)
		tt.End()
	})
}
