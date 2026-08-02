package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

// newRecording builds a T whose fatal is a recording stub so error paths can be
// exercised without aborting the enclosing test.
func newRecording(tt *tape.T, plugin any, dir string) (*T, *[]string) {
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
		tr, calls := newRecording(tt, reportPlugin{}, dir)
		tr.Report("bad", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReportZeroPlaces(t *testing.T) {
	dir := writeDir(t, map[string]string{"clean.go": cleanSrc})
	tape.Test(t, "test: Report zero places", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin{}, dir)
		tr.Report("clean", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoReport parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin{}, dir)
		tr.NoReport("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoReportHasPlaces(t *testing.T) {
	// plugin reports on this fixture, so NoReport must fail
	dir := writeDir(t, map[string]string{"match.go": matchSrc})
	tape.Test(t, "test: NoReport with places", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin{}, dir)
		tr.NoReport("match")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: Transform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin{}, dir)
		tr.Transform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestNoTransformParseError(t *testing.T) {
	dir := writeDir(t, map[string]string{"bad.go": badSrc})
	tape.Test(t, "test: NoTransform parse error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin{}, dir)
		tr.NoTransform("bad")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestReadMissingFile(t *testing.T) {
	dir := writeDir(t, map[string]string{})
	tape.Test(t, "test: read missing file", func(tt *tape.T) {
		tr, calls := newRecording(tt, reportPlugin{}, dir)
		tr.Report("nonexistent", "found it")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}

func TestTransformUpdateWriteError(t *testing.T) {
	dir := writeDir(t, map[string]string{"replace.go": replaceSrc})
	t.Setenv("UPDATE", "1")
	tape.Test(t, "test: Transform UPDATE write error", func(tt *tape.T) {
		tr, calls := newRecording(tt, replacePlugin{}, dir)
		tr.writeFile = func(_ string, _ []byte, _ os.FileMode) error {
			return os.ErrPermission
		}
		tr.Transform("replace")
		assertFatal(tt.TB(), *calls)
		tt.End()
	})
}
