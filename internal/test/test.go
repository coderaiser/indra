package test

import (
	"os"
	"path/filepath"
	"testing"

	loader "coderaiser/indra/engine-loader"
	processor "coderaiser/indra/engine-processor"
	runner "coderaiser/indra/engine-runner"
	"coderaiser/indra/internal/plugins"
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

// CreateTest returns a typed test runner bound to the plugin at pkgPath.
// Call once at package level: CreateTest("pkg/path", runtime.Caller).
// The caller function is invoked with 0 to locate the calling file; the
// fixture dir is derived as <file>/fixture.
func CreateTest(pkgPath string, caller func(int) (uintptr, string, int, bool)) func(*testing.T, string, func(*T)) {
	_, file, _, _ := caller(0)
	dir := filepath.Join(filepath.Dir(file), "fixture")
	return tape.Extend(func(base *tape.T) *T {
		return New(base, loadPlugin(pkgPath), dir)
	})
}

// loadPlugin resolves the plugin for pkgPath from the static plugin index.
func loadPlugin(pkgPath string) []runner.PluginItem {
	for _, pf := range plugins.All {
		if pf.Path != pkgPath {
			continue
		}
		kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
		return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
	}
	panic("internal/test: unknown plugin " + pkgPath)
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

