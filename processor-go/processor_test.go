package processor_go

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	loader "coderaiser/indra/engine-loader"
	runner "coderaiser/indra/engine-runner"
	"coderaiser/indra/types"
)

func pluginItems() []runner.PluginItem {
	kinds := loader.Load([]loader.PluginFuncs{{
		Name:    "eq",
		Report:  func() string { return "use DeepEqual" },
		Match:   func() types.Matcher { return types.Matcher{"t.Equal(__a, __b)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"} },
	}}, nil, loader.Config{})
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestProcessFileReturnsPlaces(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessFile(path, pluginItems(), false)
	if err != nil {
		t.Fatalf("ProcessFile: %v", err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestProcessFileFixRewrites(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	_, err := ProcessFile(path, pluginItems(), true)
	if err != nil {
		t.Fatalf("ProcessFile: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(data), "t.DeepEqual(a, b)") {
		t.Fatalf("expected rewritten file:\n%s", data)
	}
}

func TestProcessFileSkipsBuildIgnore(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "//go:build ignore\n\npackage p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessFile(path, pluginItems(), true)
	if err != nil {
		t.Fatalf("ProcessFile: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("expected no places for build-ignored file, got %d", len(places))
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "t.Equal(a, b)") {
		t.Fatal("build-ignored file must be left untouched")
	}
}

func TestProcessFileReadError(t *testing.T) {
	_, err := ProcessFile(filepath.Join(t.TempDir(), "missing.go"), pluginItems(), false)
	if err == nil {
		t.Fatal("expected read error")
	}
}

func TestProcessFileWriteError(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	if err := os.Chmod(path, 0444); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	defer os.Chmod(path, 0644)
	_, err := ProcessFile(path, pluginItems(), true)
	if err == nil {
		t.Fatal("expected write error")
	}
}

func TestProcessFileParseError(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\nfunc (\n")
	_, err := ProcessFile(path, pluginItems(), false)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestProcessDir(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	write(t, dir, "b.go", "package p\n\nfunc f() {\n\tt.Equal(c, d)\n}\n")
	write(t, dir, "c.txt", "not go")
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	write(t, sub, "d.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginItems(), false)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 2 {
		t.Fatalf("expected 2 places (non-recursive), got %d", len(places))
	}
}

func TestProcessDirReadError(t *testing.T) {
	_, err := ProcessDir(filepath.Join(t.TempDir(), "nope"), pluginItems(), false)
	if err == nil {
		t.Fatal("expected read dir error")
	}
}

func TestProcessDirLoopError(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	write(t, dir, "bad.go", "package p\nfunc (\n")
	_, err := ProcessDir(dir, pluginItems(), false)
	if err == nil {
		t.Fatal("expected loop error from invalid .go file")
	}
}
