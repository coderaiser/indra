package test

import (
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

// T wraps tape.T with plugin-level lint assertions.
type T struct {
	*tape.T
	plugins   []runner.PluginItem
	dir       string
	fatal     func(format string, args ...any)
	writeFile func(string, []byte, os.FileMode) error
}

// New constructs a T for direct use in error-path tests.
func New(tt *tape.T, plugins []runner.PluginItem, dir string) *T {
	return &T{T: tt, plugins: plugins, dir: dir, fatal: tt.TB().Fatalf, writeFile: os.WriteFile}
}

// For returns a test runner fixed to a single leaf rule whose plugin is a
// method-bearing struct (e.g. remove_skip.Plugin{}). The fixture dir is the
// caller's fixture/ directory.
//
//	func TestRemoveSkip(t *testing.T) {
//		For("remove-skip", remove_skip.Plugin{}) ...
//	}
func For(rule string, plugin any) func(*testing.T, string, func(*T)) {
	items := itemsFor(rule, plugin)
	dir := callerFixtureDir(1)
	return tape.Extend(func(base *tape.T) *T {
		return New(base, items, dir)
	})
}

// CreateTest is the exported alias for For, used by the
// convert-for-to-create-test plugin and by migrated plugin test files.
var CreateTest = For

// ForGroup returns a test runner for every rule of a group. It is used by a
// group's own test file to exercise all sub-rules at once (e.g. tape.Rules()).
func ForGroup(name string, rules []types.Rule) func(*testing.T, string, func(*T)) {
	kinds := loader.Load([]loader.PluginFuncs{{Name: name, Rules: rules}}, loader.Config{})
	items := make([]runner.PluginItem, 0, len(kinds))
	for _, k := range kinds {
		validatePlugin(k)
		items = append(items, runner.PluginItem{Rule: k.Name(), Plugin: k})
	}
	dir := callerFixtureDir(1)
	return tape.Extend(func(base *tape.T) *T {
		return New(base, items, dir)
	})
}

// callerFixtureDir returns the fixture/ directory next to the test file that
// called For/ForGroup. depth is the position of For/ForGroup's caller frame.
func callerFixtureDir(depth int) string {
	_, file, _, _ := runtime.Caller(depth + 1)
	return filepath.Join(filepath.Dir(file), "fixture")
}

// itemsFor resolves a single leaf rule into runnable PluginItems.
func itemsFor(rule string, plugin any) []runner.PluginItem {
	kinds := loader.Load([]loader.PluginFuncs{{Name: rule, Plugin: plugin}}, loader.Config{})
	validatePlugin(kinds[0])
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

// validatePlugin enforces consistency on a resolved ReplacerPlugin before it is
// used in tests: every Match entry must have a non-nil guard function and every
// Match key must also appear as a Replace key. A malformed plugin would
// otherwise pass the tester silently and fail only at fix time.
func validatePlugin(kind loader.PluginKind) {
	rp, ok := kind.(loader.ReplacerPlugin)
	if !ok {
		return
	}
	matcher := rp.Match()
	replacer := rp.Replace()
	for pattern, guard := range matcher {
		if guard == nil {
			panic("internal/test: " + kind.Name() + ": nil MatchFn for pattern " + pattern)
		}
		if _, ok := replacer[pattern]; !ok {
			panic("internal/test: " + kind.Name() + ": Match key not in Replace: " + pattern)
		}
	}
}

// Report asserts the plugin emits exactly ≥1 place whose first Message equals
// message when run against fixture file <name>.go.
func (t *T) Report(name, message string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := processor.Process(processor.Params{Src: src, Fix: false, Plugins: t.plugins})
	if err != nil {
		t.fatal("Report(%q): parse error: %v", name, err)
		return
	}
	if len(res.Places) == 0 {
		t.fatal("Report(%q): expected at least one place, got none", name)
		return
	}
	t.Equal(res.Places[0].Message, message)
}

// NoReport asserts the plugin emits no places for fixture <name>.go.
func (t *T) NoReport(name string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := processor.Process(processor.Params{Src: src, Fix: false, Plugins: t.plugins})
	if err != nil {
		t.fatal("NoReport(%q): parse error: %v", name, err)
		return
	}
	if len(res.Places) != 0 {
		t.fatal("NoReport(%q): expected no places, got %d", name, len(res.Places))
		return
	}
	t.Pass("no report")
}

// Transform runs the plugin with fix=true on fixture <name>.go and asserts
// the output matches fixture <name>-fix.go.
// Set env UPDATE=1 to regenerate the fix fixture.
func (t *T) Transform(name string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := processor.Process(processor.Params{Src: src, Fix: true, Plugins: t.plugins})
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
	t.Equal(gotStr, string(fixSrc))
}

// NoTransform asserts that fix=true leaves fixture <name>.go source unchanged.
func (t *T) NoTransform(name string) {
	t.TB().Helper()
	src := t.read(name)
	res, err := processor.Process(processor.Params{Src: src, Fix: true, Plugins: t.plugins})
	if err != nil {
		t.fatal("NoTransform(%q): parse error: %v", name, err)
		return
	}
	gotStr := string(res.Out)
	t.Equal(gotStr, string(src))
}

func (t *T) read(name string) []byte {
	t.TB().Helper()
	data, err := os.ReadFile(filepath.Join(t.dir, name+".go"))
	if err != nil {
		t.fatal("read fixture %q: %v", name, err)
		return nil
	}
	return data
}
