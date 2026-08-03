package test

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// CreateTest returns a typed test runner for the plugin at the caller's
// package path.
//
// Call once at package level, passing runtime.Caller(0) directly:
//
//	var Test = indratest.CreateTest(runtime.Caller(0))
//
// runtime.Caller(0) expands its four return values directly into the
// parameters, so no thunk or path string is needed. file resolves to the
// plugin test file, the package path is derived from go.mod, and fixture/
// is found next to the test file.
func CreateTest(_ uintptr, file string, _ int, _ bool) func(*testing.T, string, func(*T)) {
	pkgPath := derivePackagePath(file)
	dir := filepath.Join(filepath.Dir(file), "fixture")
	for _, pf := range plugins.All {
		if pf.Path != pkgPath {
			continue
		}
		kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
		return tape.Extend(func(base *tape.T) *T {
			return New(base, []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}, dir)
		})
	}
	panic("internal/test: unknown plugin " + pkgPath)
}

// modInfo caches the module name and root dir, located from go.mod.
var modInfo struct {
	name string // e.g. "coderaiser/indra"
	root string // abs path of module root dir
}

func init() {
	_, thisFile, _, _ := runtime.Caller(0)
	modInfo.root, modInfo.name = findModInfo(filepath.Dir(thisFile), func(p string) error {
		_, err := os.Stat(p)
		return err
	}, readModuleName)
}

// findModInfo walks up from startDir looking for a go.mod, returning the
// module root dir and module name. stat and moduleName are injected so the
// not-found panic path can be driven from a synthetic filesystem in tests.
func findModInfo(startDir string, stat func(string) error, moduleName func(string) string) (string, string) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, "go.mod")
		if stat(candidate) == nil {
			return dir, moduleName(candidate)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("internal/test: go.mod not found")
		}
		dir = parent
	}
}

// readModuleName returns the module directive from a go.mod file.
func readModuleName(gomod string) string {
	f, err := os.Open(gomod)
	if err != nil {
		panic("internal/test: cannot open go.mod: " + err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module ")
		}
	}
	panic("internal/test: module line not found in go.mod")
}

// derivePackagePath maps a test file path to its package import path.
func derivePackagePath(file string) string {
	rel, err := filepath.Rel(modInfo.root, filepath.Dir(file))
	if err != nil {
		panic("internal/test: cannot relativize path: " + err.Error())
	}
	return modInfo.name + "/" + filepath.ToSlash(rel)
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

