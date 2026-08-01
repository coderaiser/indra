package test

import (
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/engine"
	tape "github.com/coderaiser/go-tape"
)

// T wraps tape.T with plugin-level lint assertions.
type T struct {
	*tape.T
	plugin engine.Plugin
	dir    string
}

// New constructs a T for direct use in error-path tests.
func New(tt *tape.T, plugin engine.Plugin, dir string) *T {
	return &T{T: tt, plugin: plugin, dir: dir}
}

// Report asserts the plugin emits exactly ≥1 place whose first Message equals
// message when run against fixture file <name>.go.
func (t *T) Report(name, message string) {
	t.TB().Helper()
	src := t.read(name)
	_, places, err := engine.Indra(src, []engine.Plugin{t.plugin}, false)
	if err != nil {
		t.TB().Fatalf("Report(%q): parse error: %v", name, err)
	}
	if len(places) == 0 {
		t.TB().Fatalf("Report(%q): expected at least one place, got none", name)
	}
	t.Equal(places[0].Message, message)
}

// NoReport asserts the plugin emits no places for fixture <name>.go.
func (t *T) NoReport(name string) {
	t.TB().Helper()
	src := t.read(name)
	_, places, err := engine.Indra(src, []engine.Plugin{t.plugin}, false)
	if err != nil {
		t.TB().Fatalf("NoReport(%q): parse error: %v", name, err)
	}
	if len(places) != 0 {
		t.TB().Fatalf("NoReport(%q): expected no places, got %d", name, len(places))
	}
	t.Pass("no report")
}

// Transform runs the plugin with fix=true on fixture <name>.go and asserts
// the output matches fixture <name>-fix.go.
// Set env UPDATE=1 to regenerate the fix fixture.
func (t *T) Transform(name string) {
	t.TB().Helper()
	src := t.read(name)
	got, _, err := engine.Indra(src, []engine.Plugin{t.plugin}, true)
	if err != nil {
		t.TB().Fatalf("Transform(%q): parse error: %v", name, err)
	}
	fixPath := filepath.Join(t.dir, name+"-fix.go")
	if os.Getenv("UPDATE") == "1" {
		if werr := os.WriteFile(fixPath, got, 0644); werr != nil {
			t.TB().Fatalf("Transform(%q): write fixture: %v", name, werr)
		}
		t.Pass("fixture updated")
		return
	}
	expected := t.read(name + "-fix")
	t.Equal(string(got), string(expected))
}

// NoTransform asserts that fix=true leaves fixture <name>.go source unchanged.
func (t *T) NoTransform(name string) {
	t.TB().Helper()
	src := t.read(name)
	got, _, err := engine.Indra(src, []engine.Plugin{t.plugin}, true)
	if err != nil {
		t.TB().Fatalf("NoTransform(%q): parse error: %v", name, err)
	}
	t.Equal(string(got), string(src))
}

func (t *T) read(name string) []byte {
	t.TB().Helper()
	data, err := os.ReadFile(filepath.Join(t.dir, name+".go"))
	if err != nil {
		t.TB().Fatalf("read fixture %q: %v", name, err)
	}
	return data
}

// CreateTest returns a typed test runner bound to plugin.
// Call once at package level with the fixture dir relative to the caller file.
func CreateTest(plugin engine.Plugin, dir string) func(*testing.T, string, func(*T)) {
	return func(t *testing.T, name string, fn func(*T)) {
		tape.Test(t, name, func(tt *tape.T) {
			fn(&T{T: tt, plugin: plugin, dir: dir})
		})
	}
}
