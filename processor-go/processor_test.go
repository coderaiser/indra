package processor_go

import (
	"errors"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	loader "coderaiser/indra/engine_loader"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/internal/config"
	"coderaiser/indra/types"
)

func pluginOpts() Options {
	kinds := loader.Load([]loader.PluginFuncs{{
		Name:   "eq",
		Report: func() string { return "use DeepEqual" },
		Match: func() types.Matcher {
			return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ *ast.BlockStmt) bool { return true }}
		},
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"} },
	}}, loader.Config{})
	return Opt([]runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}, false)
}

func pluginOptsFixed() Options {
	o := pluginOpts()
	o.fix = true
	return o
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
	places, err := ProcessFile(path, pluginOpts())
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
	_, err := ProcessFile(path, pluginOptsFixed())
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
	places, err := ProcessFile(path, pluginOptsFixed())
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
	_, err := ProcessFile(filepath.Join(t.TempDir(), "missing.go"), pluginOpts())
	if err == nil {
		t.Fatal("expected read error")
	}
}

func TestProcessFileWriteError(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	_, err := ProcessFile(path, pluginOptsFixed(), WithWriteFile(func(string, []byte, os.FileMode) error {
		return errors.New("write error")
	}))
	if err == nil {
		t.Fatal("expected write error")
	}
}

func TestProcessFileParseError(t *testing.T) {
	dir := t.TempDir()
	path := write(t, dir, "a.go", "package p\nfunc (\n")
	_, err := ProcessFile(path, pluginOpts())
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestProcessDirRecursive(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	write(t, dir, "b.go", "package p\n\nfunc f() {\n\tt.Equal(c, d)\n}\n")
	write(t, dir, "c.txt", "not go")
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	write(t, sub, "d.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginOpts(), nil)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 3 {
		t.Fatalf("expected 3 places (recursive), got %d", len(places))
	}
}

func TestProcessDirReadError(t *testing.T) {
	_, err := ProcessDir(filepath.Join(t.TempDir(), "nope"), pluginOpts(), nil)
	if err == nil {
		t.Fatal("expected read dir error")
	}
}

func TestProcessDirLoopError(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	write(t, dir, "bad.go", "package p\nfunc (\n")
	_, err := ProcessDir(dir, pluginOpts(), nil)
	if err == nil {
		t.Fatal("expected loop error from invalid .go file")
	}
}

func TestProcessDirSkipsVendor(t *testing.T) {
	dir := t.TempDir()
	vendor := filepath.Join(dir, "vendor")
	os.Mkdir(vendor, 0755)
	write(t, vendor, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginOpts(), config.DefaultIgnorePatterns)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places (vendor skipped), got %d", len(places))
	}
}

func TestProcessDirSkipsHidden(t *testing.T) {
	dir := t.TempDir()
	hidden := filepath.Join(dir, ".hidden")
	os.Mkdir(hidden, 0755)
	write(t, hidden, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginOpts(), config.DefaultIgnorePatterns)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places (hidden dir skipped), got %d", len(places))
	}
}

func TestProcessDirRootHiddenDir(t *testing.T) {
	// Walking a dot-prefixed dir directly must still descend into it.
	dir := t.TempDir()
	hidden := filepath.Join(dir, ".hidden")
	os.Mkdir(hidden, 0755)
	write(t, hidden, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(hidden, pluginOpts(), nil)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place (root hidden dir descended), got %d", len(places))
	}
}

func TestProcessDirSkipsTestdata(t *testing.T) {
	dir := t.TempDir()
	td := filepath.Join(dir, "testdata")
	os.Mkdir(td, 0755)
	write(t, td, "a.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginOpts(), config.DefaultIgnorePatterns)
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places (testdata skipped), got %d", len(places))
	}
}

func TestProcessDirIgnorePattern(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a_test.go", "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	places, err := ProcessDir(dir, pluginOpts(), []string{"**/*_test.go"})
	if err != nil {
		t.Fatalf("ProcessDir: %v", err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places (ignored), got %d", len(places))
	}
}

func TestResolveArgs(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "a.go", "package p\n")
	files, dirs := ResolveArgs([]string{f, dir, "./..."})
	if len(files) != 1 || files[0] != f {
		t.Fatalf("expected 1 file, got %v", files)
	}
	if len(dirs) != 2 {
		t.Fatalf("expected 2 dirs, got %v", dirs)
	}
}

func TestResolveArgsSuffixes(t *testing.T) {
	files, dirs := ResolveArgs([]string{"pkg/...", "..."})
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %v", files)
	}
	if len(dirs) != 2 {
		t.Fatalf("expected 2 dirs, got %v", dirs)
	}
}

func TestResolveArgsBareDotDotDot(t *testing.T) {
	// "/..." trims to empty dir which must become "."
	files, dirs := ResolveArgs([]string{"/..."})
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %v", files)
	}
	if len(dirs) != 1 || dirs[0] != "." {
		t.Fatalf("expected [.] dir, got %v", dirs)
	}
}

func TestResolveArgsBareEllipsis(t *testing.T) {
	// "pkg..." matches the bare ... suffix branch
	files, dirs := ResolveArgs([]string{"pkg..."})
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %v", files)
	}
	if len(dirs) != 1 || dirs[0] != "pkg" {
		t.Fatalf("expected [pkg] dir, got %v", dirs)
	}
}

func TestCollectFiles(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n")
	write(t, dir, "b.go", "package p\n")
	all := CollectFiles(nil, []string{dir}, nil)
	if len(all) != 2 {
		t.Fatalf("expected 2 files, got %d", len(all))
	}
}

func TestCollectFilesPreservesExplicit(t *testing.T) {
	dir := t.TempDir()
	explicit := write(t, dir, "a.go", "package p\n")
	write(t, dir, "b.go", "package p\n")
	all := CollectFiles([]string{explicit}, []string{dir}, nil)
	if len(all) != 3 {
		t.Fatalf("expected 3 files (1 explicit + 2 in dir), got %d", len(all))
	}
}

func TestCollectFilesSkipsNonGoAndVendor(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n")
	write(t, dir, "note.txt", "text")
	vendor := filepath.Join(dir, "vendor")
	os.Mkdir(vendor, 0755)
	write(t, vendor, "v.go", "package p\n")
	all := CollectFiles(nil, []string{dir}, config.DefaultIgnorePatterns)
	if len(all) != 1 {
		t.Fatalf("expected 1 file (vendor and .txt skipped), got %d", len(all))
	}
}

func TestCollectFilesIgnoresPattern(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a_test.go", "package p\n")
	all := CollectFiles(nil, []string{dir}, []string{"**/*_test.go"})
	if len(all) != 0 {
		t.Fatalf("expected 0 files (ignored), got %d", len(all))
	}
}

func TestCollectFilesWalkError(t *testing.T) {
	// A non-existent dir makes WalkDir report an error inside the callback.
	all := CollectFiles(nil, []string{filepath.Join(t.TempDir(), "missing")}, nil)
	if len(all) != 0 {
		t.Fatalf("expected 0 files, got %d", len(all))
	}
}

func TestCollectFilesRootHiddenDir(t *testing.T) {
	// Walking a dot-prefixed dir directly must still descend into it.
	dir := t.TempDir()
	hidden := filepath.Join(dir, ".hidden")
	os.Mkdir(hidden, 0755)
	write(t, hidden, "a.go", "package p\n")
	all := CollectFiles(nil, []string{hidden}, nil)
	if len(all) != 1 {
		t.Fatalf("expected 1 file (root hidden dir descended), got %d", len(all))
	}
}

func TestCollectFilesNormalSubdir(t *testing.T) {
	// A plain subdirectory is descended into and its files collected.
	dir := t.TempDir()
	write(t, dir, "a.go", "package p\n")
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	write(t, sub, "b.go", "package p\n")
	all := CollectFiles(nil, []string{dir}, nil)
	if len(all) != 2 {
		t.Fatalf("expected 2 files (root + normal subdir), got %d", len(all))
	}
}
