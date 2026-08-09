package indra_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	indra "coderaiser/indra"
	engine_loader "coderaiser/indra/engine_loader"
	plugin_indra "coderaiser/indra/internal/plugin_indra"
	remove_unused_variables "coderaiser/indra/internal/plugin_remove_unused_variables"
	plugin_tape "coderaiser/indra/internal/plugin_tape"

	. "github.com/coderaiser/go-tape"
)

var testRegistry = []engine_loader.PluginFuncs{
	{Name: "tape", Rules: plugin_tape.Rules()},
	{Name: "indra", Rules: plugin_indra.Rules()},
	{Name: "remove-unused-variables", Plugin: remove_unused_variables.Plugin{}},
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if error := os.WriteFile(path, []byte(content), 0644); error != nil {
		t.Fatalf("write %s: %v", name, error)
	}
	return path
}

const matchSrc = "package p\n\nfunc F() {\n\tt.Equal(1, []int{1})\n}\n"
const cleanSrc = "package p\n\nfunc F() {}\n"
const badSrc = "package p\nfunc (\n"

func TestLintReports(t *testing.T) {
	Test(t, "lint: Lint returns no error for matching source", func(t *T) {
		_, _, error := indra.Lint(testRegistry, []byte(matchSrc), false)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: Lint returns places for matching source", func(t *T) {
		_, places, _ := indra.Lint(testRegistry, []byte(matchSrc), false)
		t.Ok(len(places) > 0)
		t.End()
	})
}

func TestLintClean(t *testing.T) {
	Test(t, "lint: Lint returns no error for clean source", func(t *T) {
		_, _, error := indra.Lint(testRegistry, []byte(cleanSrc), false)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: Lint returns no places for clean source", func(t *T) {
		_, places, _ := indra.Lint(testRegistry, []byte(cleanSrc), false)
		result := len(places)
		t.Equal(result, 0)

		t.End()
	})
}

func TestLintFix(t *testing.T) {
	Test(t, "lint: Lint fix returns no error", func(t *T) {
		_, _, error := indra.Lint(testRegistry, []byte(matchSrc), true)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: Lint fix rewrites source", func(t *T) {
		out, _, _ := indra.Lint(testRegistry, []byte(matchSrc), true)
		t.Ok(strings.Contains(string(out), "expected"))
		t.End()
	})
}

func TestLintParseError(t *testing.T) {
	Test(t, "lint: Lint returns error for invalid source", func(t *T) {
		_, _, error := indra.Lint(testRegistry, []byte(badSrc), false)
		t.Ok(error)
		t.End()
	})
}

func TestIndraVersion(t *testing.T) {
	Test(t, "lint: --version prints version", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"--version"}, &buf)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: --version writes output", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--version"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})

	Test(t, "lint: -v prints version", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"-v"}, &buf)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: -v writes output", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"-v"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})
}

func TestIndraNoFiles(t *testing.T) {
	Test(t, "lint: no files returns nil", func(t *T) {
		error := indra.Indra(testRegistry, []string{}, io.Discard)
		t.NotOk(error)
		t.End()
	})
}

func TestIndraUnknownFlag(t *testing.T) {
	Test(t, "lint: unknown flag returns non-nil error", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"--unknown-flag"}, &buf)
		t.Ok(error)

		t.End()
	})

	Test(t, "lint: unknown flag error wraps ErrInvalidOption", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"--unknown-flag"}, &buf)
		t.Ok(errors.Is(error, indra.ErrInvalidOption))
		t.End()
	})

	Test(t, "lint: unknown flag writes error message to writer", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--unknown-flag"}, &buf)
		t.Ok(strings.Contains(buf.String(), "🐊 Invalid option `--unknown-flag`."))
		t.End()
	})

	Test(t, "lint: unknown flag with no files still returns error", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"--unknown-flag"}, &buf)
		t.Ok(error)

		t.End()
	})
}

func TestIndraFormatterFlag(t *testing.T) {
	Test(t, "lint: -f json-lines sets json-lines formatter", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"-f", "json-lines", dir}, &buf)
		t.Ok(strings.Contains(buf.String(), `"name"`))
		t.End()
	})

	Test(t, "lint: --format json-lines sets json-lines formatter", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--format", "json-lines", dir}, &buf)
		t.Ok(strings.Contains(buf.String(), `"name"`))
		t.End()
	})

	Test(t, "lint: --format dump sets dump formatter", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--format", "dump", dir}, &buf)
		t.Ok(strings.Contains(buf.String(), "bad.go"))
		t.End()
	})

	Test(t, "lint: --format=dump sets dump formatter via equals form", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--format=dump", dir}, &buf)
		t.Ok(strings.Contains(buf.String(), "bad.go"))
		t.End()
	})

	Test(t, "lint: -f=json-lines sets json-lines formatter", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"-f=json-lines", dir}, &buf)
		t.Ok(strings.Contains(buf.String(), `"name"`))
		t.End()
	})
}

func TestIndraMissingFileError(t *testing.T) {
	Test(t, "lint: missing file returns error", func(t *T) {
		error := indra.Indra(testRegistry, []string{"/nonexistent/file.go"}, io.Discard)
		t.Ok(error)
		t.End()
	})
}

func TestIndraCleanFile(t *testing.T) {
	Test(t, "lint: clean file returns nil", func(t *T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "clean.go", cleanSrc)
		error := indra.Indra(testRegistry, []string{f}, io.Discard)
		t.NotOk(error)
		t.End()
	})
}

func TestIndraFailsOnIssue(t *testing.T) {
	Test(t, "lint: file with issue returns error", func(t *T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		error := indra.Indra(testRegistry, []string{f}, io.Discard)
		t.Ok(error)
		t.End()
	})
}

func TestIndraFixWrites(t *testing.T) {
	Test(t, "lint: --fix returns nil", func(t *T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		error := indra.Indra(testRegistry, []string{"--fix", f}, io.Discard)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: --fix rewrites file", func(t *T) {
		dir := t.TB().TempDir()
		f := writeFile(t.TB(), dir, "bad.go", matchSrc)
		indra.Indra(testRegistry, []string{"--fix", f}, io.Discard)
		data, _ := os.ReadFile(f)
		t.Ok(strings.Contains(string(data), "expected"))
		t.End()
	})
}

func TestIndraDumpFormatterReports(t *testing.T) {
	Test(t, "lint: Indra dir with dump formatter reports findings", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{dir}, &buf)
		t.Ok(error != nil && buf.Len() > 0)
		t.End()
	})

	Test(t, "lint: Indra dump formatter output names file", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{dir}, &buf)
		t.Ok(strings.Contains(buf.String(), "bad.go"))
		t.End()
	})
}

func TestIndraformatter_json_linesFormatter(t *testing.T) {
	Test(t, "lint: Indra dir with json-lines formatter outputs JSON", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "bad.go", matchSrc)
		t.TB().Setenv("INDRA_FORMATTER", "json-lines")
		t.TB().Setenv("CI", "")
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{dir}, &buf)
		t.Ok(strings.Contains(buf.String(), `"name"`))
		t.End()
	})
}

func TestIndraDumpFormatterClean(t *testing.T) {
	Test(t, "lint: Indra dir with clean files outputs nothing for dump", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "clean.go", cleanSrc)
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		t.TB().Setenv("CI", "")
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{dir}, &buf)
		t.Ok(error == nil && buf.Len() == 0)
		t.End()
	})
}

func TestIndraNoGoFiles(t *testing.T) {
	Test(t, "lint: Indra dir with no go files returns nil", func(t *T) {
		dir := t.TB().TempDir()
		writeFile(t.TB(), dir, "readme.txt", "hello")
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{dir}, &buf)
		t.Ok(error == nil && buf.Len() == 0)
		t.End()
	})
}

func TestIndraHelp(t *testing.T) {
	Test(t, "lint: --help returns nil", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: --help writes output", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})

	Test(t, "lint: -h returns nil", func(t *T) {
		var buf bytes.Buffer
		error := indra.Indra(testRegistry, []string{"-h"}, &buf)
		t.NotOk(error)
		t.End()
	})

	Test(t, "lint: -h writes output", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"-h"}, &buf)
		t.Ok(buf.Len() > 0)
		t.End()
	})

	Test(t, "lint: --help output contains Usage", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Match(buf.String(), "Usage")
		t.End()
	})

	Test(t, "lint: --help output contains --fix", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Match(buf.String(), "--fix")
		t.End()
	})

	Test(t, "lint: --help output contains --format", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Match(buf.String(), "--format")
		t.End()
	})

	Test(t, "lint: --help output contains --progress", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Match(buf.String(), "--progress")
		t.End()
	})

	Test(t, "lint: --help output contains --no-progress", func(t *T) {
		var buf bytes.Buffer
		indra.Indra(testRegistry, []string{"--help"}, &buf)
		t.Match(buf.String(), "--no-progress")
		t.End()
	})
}

func TestIndraPerFileMatchOverride(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".indra.toml", "[match]\n\"skip_*.go\" = { \"tape\" = \"off\" }\n")
	writeFile(t, dir, "match.go", matchSrc)
	writeFile(t, dir, "skip_match.go", matchSrc)
	t.Setenv("INDRA_FORMATTER", "dump")
	t.Setenv("CI", "")
	t.Setenv("INDRA_PROGRESS_BAR", "0")
	t.Chdir(dir)
	var buf bytes.Buffer
	error := indra.Indra(testRegistry, []string{"."}, &buf)
	out := buf.String()

	Test(t, "match: override keeps reports for non matched file", func(t *T) {
		t.Ok(error != nil && strings.Contains(out, "match.go"))
		t.End()
	})

	Test(t, "match: override silences matched file", func(t *T) {
		t.NotOk(strings.Contains(out, "skip_match.go"))

		t.End()
	})
}

func TestIndraPluginsRestriction(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".indra.toml", "plugins = [\"remove-unused-variable\"]\n")
	writeFile(t, dir, "a.go", matchSrc)
	t.Setenv("INDRA_FORMATTER", "dump")
	t.Setenv("CI", "")
	t.Setenv("INDRA_PROGRESS_BAR", "0")
	t.Chdir(dir)
	var buf bytes.Buffer
	error := indra.Indra(testRegistry, []string{"."}, &buf)

	Test(t, "plugins: restriction drops non listed rules", func(t *T) {
		t.Ok(error == nil && buf.Len() == 0)
		t.End()
	})
}

func TestIndraPluginsIncludesRule(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".indra.toml", "plugins = [\"remove-unused-variable\", \"tape\"]\n")
	writeFile(t, dir, "a.go", matchSrc)
	t.Setenv("INDRA_FORMATTER", "dump")
	t.Setenv("CI", "")
	t.Setenv("INDRA_PROGRESS_BAR", "0")
	t.Chdir(dir)
	var buf bytes.Buffer
	error := indra.Indra(testRegistry, []string{"."}, &buf)

	Test(t, "plugins: list keeps listed group active", func(t *T) {
		t.Ok(error != nil && strings.Contains(buf.String(), "a.go"))
		t.End()
	})
}
