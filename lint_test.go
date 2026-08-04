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
	if error := os.WriteFile(path, []byte(content), 0644); error != nil {
		t.Fatalf("write %s: %v", name, error)
	}
	return path
}

const matchSrc = "package p\n\nfunc f() {\n\tt.Equal(1, []int{1})\n}\n"
const cleanSrc = "package p\n\nfunc f() {}\n"
const badSrc = "package p\nfunc (\n"

func TestLintReports(t *testing.T) {
	tape.Test(t, "lint: Lint returns no error for matching source", func(t *tape.T) {
		_, _, error := indra.Lint([]byte(matchSrc), false)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: Lint returns places for matching source", func(t *tape.T) {
		_, places, _ := indra.Lint([]byte(matchSrc), false)
		t.Ok(len(places) > 0)
		t.End()
	})
}

func TestLintClean(t *testing.T) {
	tape.Test(t, "lint: Lint returns no error for clean source", func(t *tape.T) {
		_, _, error := indra.Lint([]byte(cleanSrc), false)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: Lint returns no places for clean source", func(t *tape.T) {
		_, places, _ := indra.Lint([]byte(cleanSrc), false)
		t.Equal(len(places), 0)
		t.End()
	})
}

func TestLintFix(t *testing.T) {
	tape.Test(t, "lint: Lint fix returns no error", func(t *tape.T) {
		_, _, error := indra.Lint([]byte(matchSrc), true)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: Lint fix rewrites source", func(t *tape.T) {
		out, _, _ := indra.Lint([]byte(matchSrc), true)
		t.Ok(strings.Contains(string(out), "expected"))
		t.End()
	})
}

func TestLintParseError(t *testing.T) {
	tape.Test(t, "lint: Lint returns error for invalid source", func(t *tape.T) {
		_, _, error := indra.Lint([]byte(badSrc), false)
		t.Ok(error)
		t.End()
	})
}

func TestIndraVersion(t *testing.T) {
	tape.Test(t, "lint: --version prints version", func(t *tape.T) {
		var buf bytes.Buffer
		error := indra.Indra([]string{"--version"}, &buf)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: --version writes output", func(t *tape.T) {
		var buf bytes.Buffer
		indra.Indra([]string{"--version"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})

	tape.Test(t, "lint: -v prints version", func(t *tape.T) {
		var buf bytes.Buffer
		error := indra.Indra([]string{"-v"}, &buf)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: -v writes output", func(t *tape.T) {
		var buf bytes.Buffer
		indra.Indra([]string{"-v"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})
}

func TestIndraNoFiles(t *testing.T) {
	tape.Test(t, "lint: no files returns nil", func(t *tape.T) {
		error := indra.Indra([]string{}, io.Discard)
		t.NotOk(error)
		t.End()
	})
}

func TestIndraUnknownFlag(t *testing.T) {
	tape.Test(t, "lint: unknown flag with no files returns nil", func(t *tape.T) {
		error := indra.Indra([]string{"--unknown"}, io.Discard)
		t.NotOk(error)
		t.End()
	})
}

func TestIndraMissingFileError(t *testing.T) {
	tape.Test(t, "lint: missing file returns error", func(t *tape.T) {
		error := indra.Indra([]string{"/nonexistent/file.go"}, io.Discard)
		t.Ok(error)
		t.End()
	})
}

func TestIndraCleanFile(t *testing.T) {
	tape.Test(t, "lint: clean file returns nil", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "clean.go", cleanSrc)
		error := indra.Indra([]string{f}, io.Discard)
		t.NotOk(error)
		t.End()
	})
}

func TestIndraFailsOnIssue(t *testing.T) {
	tape.Test(t, "lint: file with issue returns error", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		error := indra.Indra([]string{f}, io.Discard)
		t.Ok(error)
		t.End()
	})
}

func TestIndraFixWrites(t *testing.T) {
	tape.Test(t, "lint: --fix returns nil", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		error := indra.Indra([]string{"--fix", f}, io.Discard)
		t.NotOk(error)
		t.End()
	})

	tape.Test(t, "lint: --fix rewrites file", func(t *tape.T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		indra.Indra([]string{"--fix", f}, io.Discard)
		data, _ := os.ReadFile(f)
		t.Ok(strings.Contains(string(data), "expected"))
		t.End()
	})
}

func TestIndraDumpFormatterReports(t *testing.T) {
	tape.Test(t, "lint: Indra dir with dump formatter reports findings", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		var buf bytes.Buffer
		error := indra.Indra([]string{dir}, &buf)
		t.Ok(error != nil && buf.Len() > 0)
		t.End()
	})

	tape.Test(t, "lint: Indra dump formatter output names file", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		var buf bytes.Buffer
		indra.Indra([]string{dir}, &buf)
		t.Ok(strings.Contains(buf.String(), "bad.go"))
		t.End()
	})
}

func TestIndraJsonLinesFormatter(t *testing.T) {
	tape.Test(t, "lint: Indra dir with json-lines formatter outputs JSON", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "json-lines")
		t.TB().Setenv("CI", "")
		var buf bytes.Buffer
		indra.Indra([]string{dir}, &buf)
		t.Ok(strings.Contains(buf.String(), `"name"`))
		t.End()
	})
}

func TestIndraDumpFormatterClean(t *testing.T) {
	tape.Test(t, "lint: Indra dir with clean files outputs nothing for dump", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "clean.go", cleanSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		var buf bytes.Buffer
		error := indra.Indra([]string{dir}, &buf)
		t.Ok(error == nil && buf.Len() == 0)
		t.End()
	})
}

func TestIndraNoGoFiles(t *testing.T) {
	tape.Test(t, "lint: Indra dir with no go files returns nil", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "readme.txt", "hello")
		var buf bytes.Buffer
		error := indra.Indra([]string{dir}, &buf)
		t.Ok(error == nil && buf.Len() == 0)
		t.End()
	})
}

