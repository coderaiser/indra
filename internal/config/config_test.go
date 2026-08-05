package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/config"
	tape "github.com/coderaiser/go-tape"
)

func TestLoadMissingFileNoError(t *testing.T) {
	tape.Test(t, "config: Load returns no error for missing file", func(t *tape.T) {
		_, err := config.Load(t.TB().TempDir())
		t.Ok(err == nil)
		t.End()
	})
}

func TestLoadMissingFileEmptyPatterns(t *testing.T) {
	tape.Test(t, "config: Load returns empty patterns for missing file", func(t *tape.T) {
		cfg, _ := config.Load(t.TB().TempDir())
		t.Equal(len(cfg.Ignore.Patterns), 0)
		t.End()
	})
}

func TestLoadWithPatternsNoError(t *testing.T) {
	tape.Test(t, "config: Load returns no error for .indra.toml", func(t *tape.T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`
[ignore]
patterns = ["vendor/**", "testdata/**"]
`), 0644)
		_, err := config.Load(dir)
		t.Ok(err == nil)
		t.End()
	})
}

func TestLoadWithPatternsCount(t *testing.T) {
	tape.Test(t, "config: Load returns 2 patterns from .indra.toml", func(t *tape.T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`
[ignore]
patterns = ["vendor/**", "testdata/**"]
`), 0644)
		cfg, _ := config.Load(dir)
		t.Equal(len(cfg.Ignore.Patterns), 2)
		t.End()
	})
}

func TestLoadProgressColor(t *testing.T) {
	tape.Test(t, "config: Load parses [progress] color", func(t *tape.T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte("[progress]\ncolor = \"#ff0000\"\nminCount = 5\n"), 0644)
		cfg, _ := config.Load(dir)
		t.Equal(cfg.Progress.Color, "#ff0000")
		t.End()
	})
}

func TestLoadProgressMinCount(t *testing.T) {
	tape.Test(t, "config: Load parses [progress] minCount", func(t *tape.T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte("[progress]\ncolor = \"#ff0000\"\nminCount = 5\n"), 0644)
		cfg, _ := config.Load(dir)
		t.Equal(cfg.Progress.MinCount, 5)
		t.End()
	})
}

func TestLoadMalformed(t *testing.T) {
	tape.Test(t, "config: Load returns error for malformed toml", func(t *tape.T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`[bad`), 0644)
		_, err := config.Load(dir)
		t.Ok(err != nil)
		t.End()
	})
}

func TestIsIgnoredExact(t *testing.T) {
	tape.Test(t, "isIgnored: matches exact filename", func(t *tape.T) {
		t.Ok(config.IsIgnored([]string{"foo.go"}, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredSingleStar(t *testing.T) {
	tape.Test(t, "isIgnored: matches single star glob", func(t *tape.T) {
		t.Ok(config.IsIgnored([]string{"*.go"}, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredDoubleStarPrefix(t *testing.T) {
	tape.Test(t, "isIgnored: matches double star prefix", func(t *tape.T) {
		t.Ok(config.IsIgnored([]string{"vendor/**"}, "vendor/pkg/file.go"))
		t.End()
	})
}

func TestIsIgnoredDoubleStarSuffix(t *testing.T) {
	tape.Test(t, "isIgnored: matches double star suffix", func(t *tape.T) {
		t.Ok(config.IsIgnored([]string{"**/*_test.go"}, "internal/foo/foo_test.go"))
		t.End()
	})
}

func TestIsIgnoredNoMatch(t *testing.T) {
	tape.Test(t, "isIgnored: no match returns false", func(t *tape.T) {
		t.Ok(!config.IsIgnored([]string{"vendor/**"}, "internal/foo.go"))
		t.End()
	})
}

func TestIsIgnoredEmpty(t *testing.T) {
	tape.Test(t, "isIgnored: empty patterns returns false", func(t *tape.T) {
		t.Ok(!config.IsIgnored(nil, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredZeroSegments(t *testing.T) {
	tape.Test(t, "isIgnored: double star matches zero segments", func(t *tape.T) {
		t.Ok(config.IsIgnored([]string{"vendor/**"}, "vendor/file.go"))
		t.End()
	})
}

func TestIsIgnoredDoubleStarExhausts(t *testing.T) {
	tape.Test(t, "isIgnored: double star backtracking exhausts and returns false", func(t *tape.T) {
		t.Ok(!config.IsIgnored([]string{"**/nomatch"}, "a/b/c"))
		t.End()
	})
}
