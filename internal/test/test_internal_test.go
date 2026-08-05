package test

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	loader "coderaiser/indra/engine_loader"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/types"

	tape "github.com/coderaiser/go-tape"
)

// ── fixtures ─────────────────────────────────────────────────────────────────

const matchSrc = `package fixture

func f() {
	t.Equal(a, b)
}
`

const cleanSrc = `package fixture

func f() {}
`

const replaceSrc = `package fixture

func f() {
	t.Equal(a, b)
}
`

const badSrc = "package p\nfunc (\n"

// ── synthetic plugins ────────────────────────────────────────────────────────

// synthReplacer is a field-carrying replacer for harness tests.
type synthReplacer struct {
	report  string
	match   types.Matcher
	replace types.Replacer
}

func (s synthReplacer) Report() string          { return s.report }
func (s synthReplacer) Match() types.Matcher    { return s.match }
func (s synthReplacer) Replace() types.Replacer { return s.replace }

// synthTraverser is a minimal traverser for harness tests.
type synthTraverser struct{}

func (synthTraverser) Report(_ ast.Node) string { return "t" }
func (synthTraverser) Traverse() types.Traverser {
	return types.Traverser{"*ast.File": func(ast.Node, func(ast.Node)) {}}
}
func (synthTraverser) Fix(_ ast.Node, _ map[string]any) {}

// items builds runnable PluginItems from a synthetic replacer plugin.
func items(report string, match types.Matcher, replace types.Replacer) []runner.PluginItem {
	pf := loader.PluginFuncs{Name: "synth", Plugin: synthReplacer{report: report, match: match, replace: replace}}
	kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

// ── helpers ──────────────────────────────────────────────────────────────────

// newRecording builds a T whose fatal is a recording stub so error paths can be
// exercised without aborting the enclosing test.
func newRecording(tt *tape.T, plugins []runner.PluginItem, dir string) (*T, *[]string) {
	tr := New(tt, plugins, dir)
	calls := &[]string{}
	tr.fatal = func(format string, args ...any) {
		*calls = append(*calls, fmt.Sprintf(format, args...))
	}
	return tr, calls
}

func writeDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0644); err != nil {
			t.Fatalf("writeDir: %v", err)
		}
	}
	return dir
}

// assertFatal fails the test unless at least one fatal error was recorded.
func assertFatal(t *testing.T, calls []string) {
	t.Helper()
	if len(calls) == 0 {
		t.Error("expected a fatal error to be recorded")
	}
}

func catchPanic(fn func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	fn()
	return ""
}

// ── For / ForGroup ───────────────────────────────────────────────────────────

func TestFor(t *testing.T) {
	run := For("synth", synthReplacer{
		report:  "x",
		match:   types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }},
		replace: types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"},
	})
	run(t, "for: returns working runner", func(tt *T) {
		tt.Ok(true)
		tt.End()
	})
}

func TestForGroup(t *testing.T) {
	rules := []types.Rule{
		{Name: "a", Plugin: synthReplacer{
			report:  "x",
			match:   types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }},
			replace: types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"},
		}},
		{Name: "b", Plugin: synthTraverser{}},
	}
	run := ForGroup("g", rules)
	run(t, "forgroup: returns working runner", func(tt *T) {
		tt.Ok(true)
		tt.End()
	})
}

// ── error-path tests ─────────────────────────────────────────────────────────

func TestReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Report parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, nil), dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, nil), dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, nil), dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportHasPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport with places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, nil), dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, nil), dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.writeFile = func(_ string, _ []byte, _ os.FileMode) error {
			return os.ErrPermission
		}
		tr.Transform("replace")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateHappy(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE writes fixture", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.Transform("replace")
		if len(*calls) != 0 {
			t.Error("expected no fatal for happy update")
		}
		tt.End()
	})
}

// ── validatePlugin ───────────────────────────────────────────────────────────

func syntheticNilGuard() loader.PluginFuncs {
	return loader.PluginFuncs{Name: "nil-guard", Plugin: synthReplacer{
		report:  "x",
		match:   types.Matcher{"p": nil},
		replace: types.Replacer{"p": "q"},
	}}
}

func syntheticOrphanKey() loader.PluginFuncs {
	return loader.PluginFuncs{Name: "orphan-key", Plugin: synthReplacer{
		report:  "x",
		match:   types.Matcher{"p": func(types.Vars, *ast.BlockStmt) bool { return true }},
		replace: types.Replacer{},
	}}
}

func syntheticTraverser() loader.PluginFuncs {
	return loader.PluginFuncs{Name: "trav", Plugin: synthTraverser{}}
}

func TestValidatePluginNilGuard(t *testing.T) {
	tape.Test(t, "validatePlugin: panics on nil MatchFn", func(t *tape.T) {
		kinds := loader.Load([]loader.PluginFuncs{syntheticNilGuard()}, loader.Config{})
		msg := catchPanic(func() { validatePlugin(kinds[0]) })
		t.Ok(strings.Contains(msg, "nil MatchFn"))
		t.End()
	})
}

func TestValidatePluginOrphanKey(t *testing.T) {
	tape.Test(t, "validatePlugin: panics on orphan Match key", func(t *tape.T) {
		kinds := loader.Load([]loader.PluginFuncs{syntheticOrphanKey()}, loader.Config{})
		msg := catchPanic(func() { validatePlugin(kinds[0]) })
		t.Ok(strings.Contains(msg, "Match key not in Replace"))
		t.End()
	})
}

func TestValidatePluginTraverserNoPanic(t *testing.T) {
	tape.Test(t, "validatePlugin: no panic for traverser plugin", func(t *tape.T) {
		kinds := loader.Load([]loader.PluginFuncs{syntheticTraverser()}, loader.Config{})
		msg := catchPanic(func() { validatePlugin(kinds[0]) })
		t.Equal(msg, "")
		t.End()
	})
}
