package test_test

import (
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/engine"
	indratest "coderaiser/indra/internal/test"
	tape "github.com/coderaiser/go-tape"
)

// ── fixtures written inline ──────────────────────────────────────────────────

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

const replacedSrc = `package fixture

func f() {
	t.DeepEqual(a, b)
}
`

const badSrc = "package p\nfunc (\n"

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

// ── dir helpers ──────────────────────────────────────────────────────────────

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

// ── normal happy-path tests (using tape) ─────────────────────────────────────

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
