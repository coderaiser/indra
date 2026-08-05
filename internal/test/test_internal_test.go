package test

import (
	"go/ast"

	"fmt"
	"os"
	"runtime"
	"strings"

	"path/filepath"
	"testing"

	loader "coderaiser/indra/engine_loader"
	runner "coderaiser/indra/engine_runner"
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
	return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
}
func (reportPlugin) Replace() types.Replacer { return nil }

type replacePlugin struct{}

func (replacePlugin) Report() string { return "found it" }
func (replacePlugin) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
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
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, nil), dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, nil), dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, nil), dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportHasPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport with places", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, nil), dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, nil), dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
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
		tr, calls := newRecording(tt, items("found it", types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
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
		t.NotOk(err)

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
	file := filepath.Join(modInfo.root, "internal/plugins/remove_skip/remove_skip_test.go")
	run := CreateTest(0, file, 0, false)
	run(t, "createtest: real plugin path returns working runner", func(tt *T) {
		tt.Ok(run)

		tt.End()
	})
}

// TestEntryPathUnknownType covers the default branch of entryPath, where the
// Nested value is neither a string nor a PluginEntry.
func TestEntryPathUnknownType(t *testing.T) {
	t.Parallel()
	tape.Test(t, "entryPath: unknown type returns empty string", func(t *tape.T) {
		result := entryPath(42)
		t.Equal(result, "")
		t.End()
	})
}

// TestEntryPathString covers the string branch of entryPath.
func TestEntryPathString(t *testing.T) {
	t.Parallel()
	tape.Test(t, "entryPath: string returns itself", func(t *tape.T) {
		result := entryPath("coderaiser/indra/fake/path")
		t.Equal(result, "coderaiser/indra/fake/path")
		t.End()
	})
}

// TestEntryPathPluginEntry covers the PluginEntry branch of entryPath.
func TestEntryPathPluginEntry(t *testing.T) {
	t.Parallel()
	tape.Test(t, "entryPath: PluginEntry returns its Path", func(t *tape.T) {
		result := entryPath(types.PluginEntry{Path: "coderaiser/indra/fake/path", Enabled: true})
		t.Equal(result, "coderaiser/indra/fake/path")
		t.End()
	})
}

// TestCreateItemsFromNestedNotInLoader covers the inner return nil in
// createItemsFrom: a nested entry whose rule path matches a Rules value but
// whose rule name is not returned by loader.Load against the injected input.
func TestCreateItemsFromNestedNotInLoader(t *testing.T) {
	t.Parallel()
	tape.Test(t, "createItemsFrom: nested member not in loader output returns nil", func(t *tape.T) {
		fakeNested := types.Nested{
			"fake-rule": "coderaiser/indra/fake/path",
		}
		fakeAll := []loader.PluginFuncs{
			{Name: "fake-group", Path: "coderaiser/indra/fake/group", Rules: fakeNested},
		}
		// loader.Load with an empty input returns nothing, so the inner loop
		// finds no matching rule name and returns nil.
		result := createItemsFrom("coderaiser/indra/fake/path", fakeAll, nil)
		t.Equal(result == nil, true)
		t.End()
	})
}

// TestCreateItemsFromTopLevelLeaf covers the first-loop branch of
// createItemsFrom where the pkgPath matches a top-level leaf entry in all.
func TestCreateItemsFromTopLevelLeaf(t *testing.T) {
	t.Parallel()
	tape.Test(t, "createItemsFrom: top-level leaf returns loadItems", func(t *tape.T) {
		fakeAll := []loader.PluginFuncs{
			{
				Name:    "fake-leaf",
				Path:    "coderaiser/indra/fake/leaf",
				Report:  func() string { return "fake" },
				Match:   func() types.Matcher { return types.Matcher{"f()": func(v types.Vars) bool { return true }} },
				Replace: func() types.Replacer { return types.Replacer{"f()": "g()"} },
			},
		}
		result := createItemsFrom("coderaiser/indra/fake/leaf", fakeAll, fakeAll)
		t.Ok(len(result) == 1 && result[0].Rule == "fake-leaf")
		t.End()
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
	if len(items) != 9 {
		t.Fatalf("expected 9 tape sub-rules, got %d", len(items))
	}
	got := map[string]bool{}
	for _, it := range items {
		got[it.Rule] = true
	}
	for _, rule := range []string{
		"tape/remove-skip",
		"tape/add-t-end",
		"tape/convert-equal-to-deep-equal",
		"tape/convert-equal-to-not-ok",
		"tape/convert-ok-to-not-ok",
		"tape/remove-useless-condition",
		"tape/extract-result-from-assertion",
		"tape/remove-useless-prefix",
	} {
		if !got[rule] {
			t.Errorf("missing rule %q", rule)
		}
	}
}

// TestCreateItemsNestedMember covers resolving a plugin registered only inside
// a nested group (remove-skip lives only in tape.Rules, not in All).
func TestCreateItemsNestedMember(t *testing.T) {
	items := createItems("coderaiser/indra/internal/plugins/remove_skip")
	if len(items) != 1 {
		t.Fatalf("expected single item for nested member, got %d", len(items))
	}
	if items[0].Rule != "tape/remove-skip" {
		t.Fatalf("expected tape/remove-skip item, got %q", items[0].Rule)
	}
}

// TestCreateItemsUnknown covers the unknown-plugin return of nil.
func TestCreateItemsUnknown(t *testing.T) {
	if items := createItems("coderaiser/indra/nonexistent"); items != nil {
		t.Fatalf("expected nil for unknown plugin, got %v", items)
	}
}

// TestLoadItemsLeaf covers the leaf-plugin branch of loadItems for a plugin
// registered as a top-level leaf (remove-useless-match is in All directly).
func TestLoadItemsLeaf(t *testing.T) {
	var pf loader.PluginFuncs
	for _, p := range plugins.All {
		if p.Path == "coderaiser/indra/internal/plugin_indra/remove_useless_match" {
			pf = p
		}
	}
	if pf.Rules != nil {
		t.Fatal("expected remove-useless-match to be a leaf plugin")
	}
	items := loadItems(pf)
	if len(items) != 1 || items[0].Rule != "remove-useless-match" {
		t.Fatalf("expected single remove-useless-match item, got %v", items)
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

func catchPanic(fn func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	fn()
	return ""
}

func syntheticNilGuard() loader.PluginFuncs {
	return loader.PluginFuncs{
		Name:    "nil-guard",
		Report:  func() string { return "x" },
		Match:   func() types.Matcher { return types.Matcher{"p": nil} },
		Replace: func() types.Replacer { return types.Replacer{"p": "q"} },
	}
}

func syntheticOrphanKey() loader.PluginFuncs {
	return loader.PluginFuncs{
		Name:    "orphan-key",
		Report:  func() string { return "x" },
		Match:   func() types.Matcher { return types.Matcher{"p": func(types.Vars) bool { return true }} },
		Replace: func() types.Replacer { return types.Replacer{} },
	}
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
		pf := loader.PluginFuncs{
			Name:     "trav",
			Report:   func(_ ast.Node) string { return "x" },
			Traverse: func() types.Traverser { return types.Traverser{"*ast.File": func(ast.Node, func(ast.Node)) {}} },
			Fix:      func(ast.Node, map[string]any) {},
		}
		kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
		msg := catchPanic(func() { validatePlugin(kinds[0]) })
		t.Equal(msg, "")
		t.End()
	})
}
