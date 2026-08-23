// Package test provides plugin acceptance-test helpers that any linter can
// drive through a types.Lint function. The engine itself stays external:
// indra binds its own engine via the internal/test shim, and go-lint binds
// its own engine the same way.
//
//	var CreateTest = indratest.CreateTest(indraLint)
//
//	var CreateTest = indratest.CreateTest(golint.Lint)
package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	loader "coderaiser/indra/engine_loader"
	"coderaiser/indra/types"

	tape "github.com/coderaiser/go-tape"
)

// PluginArg pairs a rule name with its plugin payload inside the []any slice
// handed to a types.Lint implementation, so the linter can resolve the fully
// qualified rule name when it builds runnable plugin items.
type PluginArg struct {
	Rule   string
	Plugin any
	// Config optionally carries loader rule state (options) applied when
	// resolving this plugin; nil means default (enabled, no options).
	Config loader.Config
}

// T wraps tape.T with plugin-level lint assertions.
type T struct {
	*tape.T
	lint         types.Lint
	plugins      []any
	dir          string
	fatal        func(format string, args ...any)
	writeFile    func(string, []byte, os.FileMode) error
	reportCustom func(ok bool, operatorName, output string, got, expected any)
}

// New constructs a T for direct use in error-path tests.
func New(tt *tape.T, lint types.Lint, plugins []any, dir string) *T {
	return &T{
		T:            tt,
		lint:         lint,
		plugins:      plugins,
		dir:          dir,
		fatal:        tt.TB().Fatalf,
		writeFile:    os.WriteFile,
		reportCustom: tt.ReportCustom,
	}
}

// CreateTest returns a test runner fixed to a single leaf rule whose plugin is
// a method-bearing struct (e.g. remove_skip.Plugin{}). The fixture dir is the
// caller's fixture/ directory.
//
//	func TestRemoveSkip(t *testing.T) {
//		CreateTest(lint)("remove-skip", remove_skip.Plugin{}) ...
//	}
func CreateTest(lint types.Lint) func(rule string, plugin any) func(*testing.T, string, func(*T)) {
	return func(rule string, plugin any) func(*testing.T, string, func(*T)) {
		plugins := []any{PluginArg{Rule: rule, Plugin: plugin}}
		dir := callerFixtureDir(1)
		return tape.Extend(func(base *tape.T) *T {
			return New(base, lint, plugins, dir)
		})
	}
}

// callerFixtureDir returns the fixture/ directory next to the test file that
// called CreateTest. depth is the position of CreateTest's caller frame.
func callerFixtureDir(depth int) string {
	_, file, _, _ := runtime.Caller(depth + 1)
	return filepath.Join(filepath.Dir(file), "fixture")
}

// Report asserts the plugin emits at least one place whose first Message
// equals message when run against fixture file <name>.go. Emits operator name
// "report" to match @putout/test.
func (t *T) Report(name, message string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := t.lint(src, false, t.plugins)
	if err != nil {
		t.fatal("Report(%q): parse error: %v", name, err)
		return
	}
	if len(res.Places) == 0 {
		t.fatal("Report(%q): expected at least one place, got none", name)
		return
	}
	got := res.Places[0].Message
	r := tape.BuiltinOperators.Equal(got, message)
	t.reportCustom(r.Ok, "report", r.Output, r.Result, r.Expected)
}

// NoReport asserts the plugin emits no places for fixture <name>.go.
// Emits operator name "noReport" to match @putout/test.
func (t *T) NoReport(name string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := t.lint(src, false, t.plugins)
	if err != nil {
		t.fatal("NoReport(%q): parse error: %v", name, err)
		return
	}
	if len(res.Places) != 0 {
		t.fatal("NoReport(%q): expected no places, got %d", name, len(res.Places))
		return
	}
	t.reportCustom(true, "noReport", "", nil, nil)
}

// Transform runs the plugin with fix=true on fixture <name>.go and asserts
// the output matches fixture <name>-fix.go.
// Set env UPDATE=1 to regenerate the fix fixture.
// Emits operator name "transform" to match @putout/test.
func (t *T) Transform(name string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := t.lint(src, true, t.plugins)
	if err != nil {
		t.fatal("Transform(%q): parse error: %v", name, err)
		return
	}
	fixPath := filepath.Join(t.dir, name+"-fix.go")
	if os.Getenv("UPDATE") == "1" {
		if werr := t.writeFile(fixPath, res.Out, 0644); werr != nil {
			t.fatal("Transform(%q): write fixture: %v", name, werr)
			return
		}
		t.Pass("fixture updated")
		return
	}
	fixSrc := t.read(name + "-fix")
	gotStr := string(res.Out)
	r := tape.BuiltinOperators.Equal(gotStr, string(fixSrc))
	t.reportCustom(r.Ok, "transform", r.Output, r.Result, r.Expected)
}

// NoTransform asserts that fix=true leaves fixture <name>.go source unchanged.
// Set env UPDATE=1 to delete a stale <name>-fix.go fixture instead.
// Emits operator name "noTransform" to match @putout/test.
func (t *T) NoTransform(name string) {
	t.TB().Helper()
	if os.Getenv("UPDATE") == "1" {
		fixPath := filepath.Join(t.dir, name+"-fix.go")
		if err := os.Remove(fixPath); err != nil && !os.IsNotExist(err) {
			t.fatal("NoTransform(%q): remove fix fixture: %v", name, err)
			return
		}
		t.Pass("fix fixture removed")
		return
	}
	src := t.read(name)
	res, err := t.lint(src, true, t.plugins)
	if err != nil {
		t.fatal("NoTransform(%q): parse error: %v", name, err)
		return
	}
	gotStr := string(res.Out)
	r := tape.BuiltinOperators.Equal(gotStr, string(src))
	t.reportCustom(r.Ok, "noTransform", r.Output, r.Result, r.Expected)
}

func (t *T) read(name string) []byte {
	t.TB().Helper()
	path := filepath.Join(t.dir, name+".go")
	if _, err := os.Stat(path); err != nil {
		// Non-Go fixtures keep their real extension (e.g. .json).
		matches, _ := filepath.Glob(filepath.Join(t.dir, name) + ".*")
		if len(matches) > 0 {
			path = matches[0]
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.fatal("read fixture %q: %v", name, err)
		return nil
	}
	return data
}
