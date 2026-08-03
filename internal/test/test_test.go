package test_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	loader "coderaiser/indra/engine-loader"
	runner "coderaiser/indra/engine-runner"
	indratest "coderaiser/indra/internal/test"
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

// items builds runnable PluginItems from synthetic plugin funcs.
func items(report string, match types.Matcher, replace types.Replacer) []runner.PluginItem {
	pf := loader.PluginFuncs{Name: "synth", Report: func() string { return report }, Match: func() types.Matcher { return match }, Replace: func() types.Replacer { return replace }}
	kinds := loader.Load([]loader.PluginFuncs{pf}, loader.Config{})
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

// testRunner wraps tape.Extend around indratest.New for synthetic plugins.
func testRunner(plugins []runner.PluginItem, dir string) func(*testing.T, string, func(*indratest.T)) {
	return tape.Extend(func(base *tape.T) *indratest.T {
		return indratest.New(base, plugins, dir)
	})
}

// ── dir helper ───────────────────────────────────────────────────────────────

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

// ── happy-path tests ─────────────────────────────────────────────────────────

func TestReport(t *testing.T) {
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	Test := testRunner(items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
	Test(t, "test: Report correct message", func(t *indratest.T) {
		t.Report("match", "found it")
		t.End()
	})
}

func TestNoReport(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	Test := testRunner(items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
	Test(t, "test: NoReport clean fixture", func(t *indratest.T) {
		t.NoReport("clean")
		t.End()
	})
}

func TestTransform(t *testing.T) {
	dir := writeDir(t, map[string]string{
		"replace.go":     replaceSrc,
		"replace-fix.go": replacedSrc,
	})
	Test := testRunner(items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
	Test(t, "test: Transform matches fix fixture", func(t *indratest.T) {
		t.Transform("replace")
		t.End()
	})
}

func TestTransformUpdate(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	Test := testRunner(items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}), dir)
	Test(t, "test: Transform UPDATE=1 writes fix fixture", func(t *indratest.T) {
		t.Transform("replace")
		t.End()
	})
	data, err := os.ReadFile(filepath.Join(dir, "replace-fix.go"))
	if err != nil {
		t.Fatalf("expected fix file to be written: %v", err)
	}
	if string(data) != replacedSrc {
		t.Fatalf("fix file content wrong:\ngot:  %q\nwant: %q", data, replacedSrc)
	}
}

func TestNoTransform(t *testing.T) {
	// report-only plugin → no rewrite → src must be unchanged
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	Test := testRunner(items("found it", types.Matcher{"t.Equal(__a, __b)": nil}, nil), dir)
	Test(t, "test: NoTransform unchanged fixture", func(t *indratest.T) {
		t.NoTransform("replace")
		t.End()
	})
}

// TestCreateTestRealPlugin covers CreateTest + loadPlugin on a real plugin
// from the static index, exercising the New + loadPlugin happy path. The
// fixture lives in the package's fixture/ dir (derived from runtime.Caller).
func TestCreateTestRealPlugin(t *testing.T) {
	Test := indratest.CreateTest("coderaiser/indra/internal/plugins/remove-skip", runtime.Caller)
	Test(t, "test: CreateTest real plugin reports", func(t *indratest.T) {
		t.Report("skip", "remove Test.Skip call")
		t.End()
	})
}

