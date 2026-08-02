package lint

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRunParseError(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad.go")
	if err := os.WriteFile(badFile, []byte("package p\n\nfunc (\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	failed := Run([]string{badFile}, &buf)
	if !failed {
		t.Error("expected failed=true for parse error")
	}
	if buf.Len() == 0 {
		t.Error("expected stderr output for parse error")
	}
}

func TestRunCleanFile(t *testing.T) {
	tmpDir := t.TempDir()
	cleanFile := filepath.Join(tmpDir, "clean.go")
	if err := os.WriteFile(cleanFile, []byte("package p\n\nfunc f() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	failed := Run([]string{cleanFile}, &buf)
	if failed {
		t.Error("expected failed=false for clean file")
	}
	if buf.Len() != 0 {
		t.Error("expected no output for clean file")
	}
}

func TestFixRuleFires(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "fix_test.go")
	src := `package p

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(1, []int{1, 2, 3})
}
`
	if err := os.WriteFile(srcFile, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	// writeFile spy that captures writes
	var writtenFile string
	var writtenContent []byte
	spyWriteFile := func(path string, content []byte, mode os.FileMode) error {
		writtenFile = path
		writtenContent = content
		return nil
	}

	failed := run([]string{srcFile}, &buf, true, spyWriteFile)
	if writtenFile == "" {
		t.Error("expected writeFile to be called")
	}
	if writtenContent == nil {
		t.Error("expected written content")
	}
	_ = failed
}

func TestRunReportsResults(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "report_test.go")
	src := `package p

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(1, []int{1})
}
`
	if err := os.WriteFile(srcFile, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	failed := Run([]string{srcFile}, &buf)
	if !failed {
		t.Error("expected failed=true for file with issues")
	}
	if buf.Len() == 0 {
		t.Error("expected output for file with issues")
	}
}

func TestWriteFileError(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "write_err_test.go")
	src := `package p

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(1, []int{1, 2, 3})
}
`
	if err := os.WriteFile(srcFile, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	errWriteFile := func(path string, content []byte, mode os.FileMode) error {
		return io.ErrClosedPipe
	}

	_ = run([]string{srcFile}, &buf, true, errWriteFile)
	if buf.Len() == 0 {
		t.Error("expected error output for writeFile failure")
	}
}

func TestRunNoFiles(t *testing.T) {
	var buf bytes.Buffer
	failed := Run([]string{}, &buf)
	if failed {
		t.Error("expected failed=false for no files")
	}
}

func TestFixNoFiles(t *testing.T) {
	var buf bytes.Buffer
	failed := Fix([]string{}, &buf)
	if failed {
		t.Error("expected failed=false for no files")
	}
}
func TestRunReadError(t *testing.T) {
	var buf bytes.Buffer
	failed := Run([]string{"/nonexistent/path/file.go"}, &buf)
	if !failed {
		t.Error("expected failed=true for unreadable file")
	}
	if buf.Len() == 0 {
		t.Error("expected error output for unreadable file")
	}
}
