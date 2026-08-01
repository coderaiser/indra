package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/engine"
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

// newRecording builds a T whose fatal is a recording stub so error paths can be
// exercised without aborting the enclosing test.
func newRecording(tt *tape.T, plugin engine.Plugin, dir string) (*T, *[]string) {
	tr := New(tt, plugin, dir)
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
		tr, calls := newRecording(tt, reportPlugin(), dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin(), dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin(), dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportHasPlaces(t *testing.T) {
	// plugin reports on this fixture, so NoReport must fail
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport with places", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin(), dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin(), dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin(), dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin(), dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Logf("skip: running as root, chmod restrictions do not apply")
		return
	}
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	// r-x: allow reading fixtures but prevent writing the fix fixture.
	if err := os.Chmod(dir, 0555); err != nil {
		t.Logf("skip: cannot chmod: %v", err)
		return
	}
	t.Cleanup(func() { os.Chmod(dir, 0755) })
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin(), dir)
		tr.Transform("replace")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}
