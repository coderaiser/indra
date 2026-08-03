package indra_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	indra "coderaiser/indra"
	tape "github.com/coderaiser/go-tape"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestLintReports(t *testing.T) {
	tape.Test(t, "lint: Lint returns a place for a matching assertion", func(t *tape.T) {
		src := []byte("package p\n\nfunc f() {\n\tt.Equal(1, []int{1})\n}\n")
		_, places, err := indra.Lint(src, false)
		t.Ok(err == nil && len(places) > 0)
		t.End()
	})
}

func TestLintClean(t *testing.T) {
	tape.Test(t, "lint: Lint returns no places for clean source", func(t *tape.T) {
		src := []byte("package p\n\nfunc f() {}\n")
		_, places, err := indra.Lint(src, false)
		t.Ok(err == nil && len(places) == 0)
		t.End()
	})
}

func TestLintFix(t *testing.T) {
	tape.Test(t, "lint: Lint rewrites with fix=true", func(t *tape.T) {
		src := []byte("package p\n\nfunc f() {\n\tt.Equal(1, []int{1})\n}\n")
		out, _, err := indra.Lint(src, true)
		t.Ok(err == nil && strings.Contains(string(out), "expected"))
		t.End()
	})
}

func TestLintParseError(t *testing.T) {
	tape.Test(t, "lint: Lint returns error for invalid source", func(t *tape.T) {
		_, _, err := indra.Lint([]byte("package p\nfunc (\n"), false)
		t.Ok(err != nil)
		t.End()
	})
}

func TestIndraVersion(t *testing.T) {
	tape.Test(t, "lint: --version prints version", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.Indra([]string{"--version"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})

	tape.Test(t, "lint: -v prints version", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.Indra([]string{"-v"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})
}

func TestIndraNoFiles(t *testing.T) {
	tape.Test(t, "lint: no files returns nil", func(t *tape.T) {
		err := indra.Indra([]string{}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestIndraUnknownFlag(t *testing.T) {
	tape.Test(t, "lint: unknown flag filtered out, no files returns nil", func(t *tape.T) {
		err := indra.Indra([]string{"--unknown"}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestIndraMissingFileError(t *testing.T) {
	tape.Test(t, "lint: missing file fails with error", func(t *tape.T) {
		err := indra.Indra([]string{"/nonexistent/file.go"}, io.Discard)
		t.Ok(err != nil)
		t.End()
	})
}

func TestIndraCleanFile(t *testing.T) {
	tape.Test(t, "lint: clean file returns nil", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "clean.go", "package p\n\nfunc f() {}\n")
		err := indra.Indra([]string{f}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestIndraFailsOnIssue(t *testing.T) {
	tape.Test(t, "lint: file with issue fails", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", "package p\n\nfunc f() {\n\tt.Equal(1, []int{1})\n}\n")
		err := indra.Indra([]string{f}, io.Discard)
		t.Ok(err != nil)
		t.End()
	})
}

func TestIndraFixWrites(t *testing.T) {
	tape.Test(t, "lint: --fix rewrites the file", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", "package p\n\nfunc f() {\n\tt.Equal(1, []int{1})\n}\n")
		err := indra.Indra([]string{"--fix", f}, io.Discard)
		data, _ := os.ReadFile(f)
		t.Ok(err == nil && strings.Contains(string(data), "expected"))
		t.End()
	})
}
