package indra_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	indra "coderaiser/indra"
	tape "github.com/coderaiser/go-tape"
)

func TestRunLintVersion(t *testing.T) {
	tape.Test(t, "lint: --version prints version", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.RunLint([]string{"--version"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})

	tape.Test(t, "lint: -v prints version", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.RunLint([]string{"-v"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})
}

func TestRunLintHelp(t *testing.T) {
	tape.Test(t, "lint: --help prints usage", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.RunLint([]string{"--help"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})

	tape.Test(t, "lint: -h prints usage", func(t *tape.T) {
		var buf bytes.Buffer
		err := indra.RunLint([]string{"-h"}, &buf)
		t.Ok(err == nil && buf.Len() > 0)
		t.End()
	})
}

func TestRunLintNoFiles(t *testing.T) {
	tape.Test(t, "lint: no files returns nil", func(t *tape.T) {
		err := indra.RunLint([]string{}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestRunLintUnknownFlag(t *testing.T) {
	tape.Test(t, "lint: unknown flag is filtered out, no files returns nil", func(t *tape.T) {
		err := indra.RunLint([]string{"--unknown"}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestRunLintWithFile(t *testing.T) {
	tape.Test(t, "lint: passes files to lint.Run and fails on error", func(t *tape.T) {
		err := indra.RunLint([]string{"/nonexistent/file.go"}, io.Discard)
		t.Ok(err != nil)
		t.End()
	})
}

func TestRunLintSuccess(t *testing.T) {
	tape.Test(t, "lint: returns nil when lint passes", func(t *tape.T) {
		cleanSrc := "package p\n\nfunc f() {}\n"
		tmpDir := t.TB().TempDir()
		cleanFile := tmpDir + "/clean.go"
		if err := os.WriteFile(cleanFile, []byte(cleanSrc), 0644); err != nil {
			t.TB().Fatal(err)
		}
		err := indra.RunLint([]string{cleanFile}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}