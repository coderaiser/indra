package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	loader "coderaiser/indra/engine-loader"
	runner "coderaiser/indra/engine-runner"
	"coderaiser/indra/internal/plugins"
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

// ── helper plugins ───────────────────────────────────────────────────────────

type reportPlugin struct{}

func (reportPlugin) Report() string { return "found it" }
func (reportPlugin) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": nil}
}
func (reportPlugin) Replace() types.Replacer { return nil }

type replacePlugin struct{}

func (replacePlugin) Report() string { return "found it" }
func (replacePlugin) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": nil}
}
func (replacePlugin) Replace() types.Replacer {
	return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
}

// items builds runnable PluginItems from synthetic plugin funcs.
func items(report string, match types.Matcher, replace types.Replacer) []runner.PluginItem {
	pf := loader.PluginFuncs{Name: "synth", Report: func() string { return report }, Match: func() types.Matcher { return match }, Replace: func() types.Replacer { return replace }}
	kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

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

// ── error-path tests ─────────────────────────────────────────────────────────

func TestReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Report parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportHasPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport with places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
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
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.Transform("replace")
		if len(*calls) != 0 {
			t.Error("expected no fatal for happy update")
		}
		tt.End()
	})
}

// TestCreateTestUnknownPluginPanics covers the unknown-plugin panic path of
// CreateTest, which is reached for a file path that does not resolve to a
// registered plugin.
func TestCreateTestUnknownPluginPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unknown plugin path")
		}
	}()
	CreateTest(0, "/nonexistent/does-not-exist/fake_test.go", 0, false)
}

func TestDerivePackagePath(t *testing.T) {
	tape.Test(t, "derivePackagePath: resolves this package correctly", func(t *tape.T) {
		_, file, _, _ := runtime.Caller(0)
		got := derivePackagePath(file)
		t.Equal(got, "coderaiser/indra/internal/test")
		t.End()
	})
}

func TestDerivePackagePathRelError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for relative path")
		}
	}()
	derivePackagePath("relative/x_test.go")
}

func TestModInfoRoot(t *testing.T) {
	tape.Test(t, "modinfo: root dir contains go.mod", func(t *tape.T) {
		_, err := os.Stat(filepath.Join(modInfo.root, "go.mod"))
		t.Equal(err, nil)
		t.End()
	})
}

func TestReadModuleName(t *testing.T) {
	tape.Test(t, "readModuleName: returns module name", func(t *tape.T) {
		name := readModuleName(filepath.Join(modInfo.root, "go.mod"))
		t.Equal(name, "coderaiser/indra")
		t.End()
	})
}

func TestReadModuleNameOpenError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when go.mod cannot be opened")
		}
	}()
	readModuleName("/nonexistent/does-not-exist/go.mod")
}

func TestReadModuleNameModuleLineMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when module line is missing")
		}
	}()
	f := filepath.Join(t.TempDir(), "go.mod")
	if err := os.WriteFile(f, []byte("// no module line\n"), 0644); err != nil {
		t.Fatal(err)
	}
	readModuleName(f)
}

func TestCreateTestHappyPath(t *testing.T) {
	file := filepath.Join(modInfo.root, "internal/plugins/remove-skip/remove_skip_test.go")
	run := CreateTest(0, file, 0, false)
	run(t, "createtest: real plugin path returns working runner", func(tt *T) {
		tt.Ok(run != nil)
		tt.End()
	})
}

// TestLoadItemsNested covers the nested-plugin branch of loadItems: a
// grouping plugin (tape) must expand to its "group/rule" sub-plugins.
func TestLoadItemsNested(t *testing.T) {
	var pf loader.PluginFuncs
	for _, p := range plugins.All {
		if p.Path == "coderaiser/indra/internal/plugins/tape" {
			pf = p
		}
	}
	if pf.Rules == nil {
		t.Fatal("expected tape plugin to carry nested Rules")
	}
	items := loadItems(pf)
	if len(items) != 2 {
		t.Fatalf("expected 2 tape sub-rules, got %d", len(items))
	}
	got := map[string]bool{}
	for _, it := range items {
		got[it.Rule] = true
	}
	for _, rule := range []string{"tape/remove-skip", "tape/add-t-end"} {
		if !got[rule] {
			t.Errorf("missing rule %q", rule)
		}
	}
}

// TestLoadItemsLeaf covers the leaf-plugin branch of loadItems.
func TestLoadItemsLeaf(t *testing.T) {
	var pf loader.PluginFuncs
	for _, p := range plugins.All {
		if p.Path == "coderaiser/indra/internal/plugins/remove-skip" {
			pf = p
		}
	}
	if pf.Rules != nil {
		t.Fatal("expected remove-skip to be a leaf plugin")
	}
	items := loadItems(pf)
	if len(items) != 1 || items[0].Rule != "remove-skip" {
		t.Fatalf("expected single remove-skip item, got %v", items)
	}
}


func TestFindModInfoFound(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/foo\n"), 0644); err != nil {
		t.Fatal(err)
	}
	root, name := findModInfo(dir, func(p string) error {
		_, err := os.Stat(p)
		return err
	}, readModuleName)
	if root != dir {
		t.Errorf("root: got %q, want %q", root, dir)
	}
	if name != "example.com/foo" {
		t.Errorf("name: got %q, want %q", name, "example.com/foo")
	}
}

func TestFindModInfoNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when no go.mod is found walking up")
		}
	}()
	findModInfo("/tmp", func(string) error { return os.ErrNotExist }, readModuleName)
}

