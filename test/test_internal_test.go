package test

import (
	"errors"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	loader "coderaiser/indra/engine_loader"
	processor "coderaiser/indra/engine_processor"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/types"

	tape "github.com/coderaiser/go-tape"
)

// ── fixture sources ──────────────────────────────────────────────────────────

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

// replacedSrc matches go/format output produced by the engine for replaceSrc.
const replacedSrc = `package fixture

func f() {
	t.DeepEqual(a, b)

}
`

// ── synthetic lints ──────────────────────────────────────────────────────────

// noopLint does nothing and returns an empty result.
func noopLint(_ []byte, _ bool, _ []any) (types.LintResult, error) {
	return types.LintResult{}, nil
}

// errLint always returns the given error.
func errLint(err error) types.Lint {
	return func(_ []byte, _ bool, _ []any) (types.LintResult, error) {
		return types.LintResult{}, err
	}
}

// placesLint returns src unchanged along with the given places.
func placesLint(places []types.Place) types.Lint {
	return func(src []byte, _ bool, _ []any) (types.LintResult, error) {
		return types.LintResult{Out: src, Places: places}, nil
	}
}

// engineLint binds a synthetic plugin to a bare leaf rule through the real
// indra engine, so happy paths run the actual parse→run→print pipeline.
func engineLint(rule string, plugin any) types.Lint {
	return func(src []byte, fix bool, _ []any) (types.LintResult, error) {
		kinds := loader.Load([]loader.PluginFuncs{{Name: rule, Plugin: plugin}}, loader.Config{})
		result, err := processor.Process(processor.Params{
			Src:     src,
			Fix:     fix,
			Plugins: []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}},
		})
		return types.LintResult{Out: result.Out, Places: result.Places}, err
	}
}

// ── synthetic plugin ─────────────────────────────────────────────────────────

type synthReplacer struct {
	report  string
	match   types.Matcher
	replace types.Replacer
}

func (s synthReplacer) Report() string          { return s.report }
func (s synthReplacer) Match() types.Matcher    { return s.match }
func (s synthReplacer) Replace() types.Replacer { return s.replace }

// ── helpers ──────────────────────────────────────────────────────────────────

// newRecording builds a T whose fatal is a recording stub so error paths can be
// exercised without aborting the enclosing test.
func newRecording(tt *tape.T, lint types.Lint, plugins []any, dir string) (*T, *[]string) {
	tr := New(tt, lint, plugins, dir)
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

func assertFatal(t *testing.T, calls []string) {
	t.Helper()
	if len(calls) == 0 {
		t.Error("expected a fatal error to be recorded")
	}
}

// ── CreateTest / ForGroup ────────────────────────────────────────────────────

func TestCreateTest(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	expected := filepath.Join(filepath.Dir(file), "fixture")
	var got string
	run := CreateTest(noopLint)("synth", synthReplacer{report: "x"})
	run(t, "test: CreateTest resolves caller fixture dir", func(tt *T) {
		got = tt.dir
		tt.End()
	})
	tape.Test(t, "test: CreateTest uses caller fixture dir", func(tt *tape.T) {
		tt.Equal(got, expected)
		tt.End()
	})
}

func TestForGroup(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	expected := filepath.Join(filepath.Dir(file), "fixture")
	rules := []types.Rule{
		{Name: "a", Plugin: synthReplacer{report: "x"}},
	}
	var got string
	run := ForGroup(noopLint)("g", rules)
	run(t, "test: ForGroup resolves caller fixture dir", func(tt *T) {
		got = tt.dir
		result := len(tt.plugins)
		tt.Equal(result, 1)

		tt.End()
	})
	tape.Test(t, "test: ForGroup uses caller fixture dir", func(tt *tape.T) {
		tt.Equal(got, expected)
		tt.End()
	})
}

// ── happy-path tests ─────────────────────────────────────────────────────────

func TestReport(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	replacer := synthReplacer{report: "found it", match: types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}}
	tape.Test(t, "test: Report correct message match", func(tt *tape.T) {
		tr := New(tt, engineLint("synth", replacer), []any{}, dir)
		tr.Report("match", "found it")
		tt.End()
	})
}

func TestNoReport(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	replacer := synthReplacer{report: "found it", match: types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}}
	tape.Test(t, "test: NoReport clean fixture", func(tt *tape.T) {
		tr := New(tt, engineLint("synth", replacer), []any{}, dir)
		tr.NoReport("clean")
		tt.End()
	})
}

func TestTransform(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"replace.go":     replaceSrc,
		"replace-fix.go": replacedSrc,
	})
	replacer := synthReplacer{report: "found it", match: types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}, replace: types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}}
	tape.Test(t, "test: Transform matches fix fixture replace", func(tt *tape.T) {
		tr := New(tt, engineLint("synth", replacer), []any{}, dir)
		tr.Transform("replace")
		tt.End()
	})
}

func TestNoTransform(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	replacer := synthReplacer{report: "found it", match: types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}}
	tape.Test(t, "test: NoTransform unchanged fixture replace", func(tt *tape.T) {
		tr := New(tt, engineLint("synth", replacer), []any{}, dir)
		tr.NoTransform("replace")
		tt.End()
	})
}

// ── error-path tests ─────────────────────────────────────────────────────────

func TestReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": "package p\nfunc (\n"})
	tape.Test(t, "test: Report parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, errLint(errors.New("boom")), []any{}, dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, placesLint(nil), []any{}, dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": "package p\nfunc (\n"})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, errLint(errors.New("boom")), []any{}, dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportUnexpectedPlace(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport unexpected place", func(tt *tape.T) {
		tr, calls := newRecording(tt, placesLint([]types.Place{{Message: "x"}}), []any{}, dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": "package p\nfunc (\n"})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, errLint(errors.New("boom")), []any{}, dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": "package p\nfunc (\n"})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, errLint(errors.New("boom")), []any{}, dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, noopLint, []any{}, dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, placesLint(nil), []any{}, dir)
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
		tr, calls := newRecording(tt, placesLint(nil), []any{}, dir)
		tr.Transform("replace")
		if len(*calls) != 0 {
			t.Error("expected no fatal for happy update")
		}
		tt.End()
	})
}
