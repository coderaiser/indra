package indra_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	indra "coderaiser/indra"
	tape "github.com/coderaiser/go-tape"
)

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
	tape.Test(t, "lint: unknown flag is filtered out, no files returns nil", func(t *tape.T) {
		err := indra.Indra([]string{"--unknown"}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}

func TestIndraWithFile(t *testing.T) {
	tape.Test(t, "lint: passes files to lint.Run and fails on error", func(t *tape.T) {
		err := indra.Indra([]string{"/nonexistent/file.go"}, io.Discard)
		t.Ok(err != nil)
		t.End()
	})
}

func TestIndraSuccess(t *testing.T) {
	tape.Test(t, "lint: returns nil when lint passes", func(t *tape.T) {
		cleanSrc := "package p\n\nfunc f() {}\n"
		tmpDir := t.TB().TempDir()
		cleanFile := tmpDir + "/clean.go"
		if err := os.WriteFile(cleanFile, []byte(cleanSrc), 0644); err != nil {
			t.TB().Fatal(err)
		}
		err := indra.Indra([]string{cleanFile}, io.Discard)
		t.Ok(err == nil)
		t.End()
	})
}