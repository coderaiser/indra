package test_test

import (
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/engine"
	indratest "coderaiser/indra/internal/test"
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

// ── helper plugins ───────────────────────────────────────────────────────────

func reportPlugin() engine.Plugin {
	return engine.Plugin{
		Name:   "test-report",
		Report: func() string { return "found it" },
		Match: func() map[string]engine.MatchFn {
			return map[string]engine.MatchFn{"t.Equal(__a, __b)": nil}
		},
	}
}

func replacePlugin() engine.Plugin {
	return engine.Plugin{
		Name:   "test-replace",
		Report: func() string { return "found it" },
		Match: func() map[string]engine.MatchFn {
			return map[string]engine.MatchFn{"t.Equal(__a, __b)": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
		},
	}
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
	Test := indratest.CreateTest(reportPlugin(), dir)
	Test(t, "test: Report correct message", func(t *indratest.T) {
		t.Report("match", "found it")
		t.End()
	})
}

func TestNoReport(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	Test := indratest.CreateTest(reportPlugin(), dir)
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
	Test := indratest.CreateTest(replacePlugin(), dir)
	Test(t, "test: Transform matches fix fixture", func(t *indratest.T) {
		t.Transform("replace")
		t.End()
	})
}

func TestTransformUpdate(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	Test := indratest.CreateTest(replacePlugin(), dir)
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
	Test := indratest.CreateTest(reportPlugin(), dir)
	Test(t, "test: NoTransform unchanged fixture", func(t *indratest.T) {
		t.NoTransform("replace")
		t.End()
	})
}
