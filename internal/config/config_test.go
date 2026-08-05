package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"coderaiser/indra/internal/config"
	. "github.com/coderaiser/go-tape"
)

func TestLoadMissingFileNoError(t *testing.T) {
	Test(t, "config: Load returns no error for missing file", func(t *T) {
		_, err := config.Load(t.TB().TempDir())
		t.NotOk(err)

		t.End()
	})
}

func TestLoadMissingFileEmptyPatterns(t *testing.T) {
	Test(t, "config: Load returns empty patterns for missing file", func(t *T) {
		cfg, _ := config.Load(t.TB().TempDir())
		result := len(cfg.Ignore.Patterns)
		t.Equal(result, 0)

		t.End()
	})
}

func TestLoadWithPatternsNoError(t *testing.T) {
	Test(t, "config: Load returns no error for .indra.toml", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`
[ignore]
patterns = ["vendor/**", "testdata/**"]
`), 0644)
		_, err := config.Load(dir)
		t.NotOk(err)

		t.End()
	})
}

func TestLoadWithPatternsCount(t *testing.T) {
	Test(t, "config: Load returns 2 patterns from .indra.toml", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`
[ignore]
patterns = ["vendor/**", "testdata/**"]
`), 0644)
		cfg, _ := config.Load(dir)
		result := len(cfg.Ignore.Patterns)
		t.Equal(result, 2)

		t.End()
	})
}

func TestLoadProgressColor(t *testing.T) {
	Test(t, "config: Load parses [progress] color", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte("[progress]\ncolor = \"#ff0000\"\nminCount = 5\n"), 0644)
		cfg, _ := config.Load(dir)
		t.Equal(cfg.Progress.Color, "#ff0000")
		t.End()
	})
}

func TestLoadProgressMinCount(t *testing.T) {
	Test(t, "config: Load parses [progress] minCount", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte("[progress]\ncolor = \"#ff0000\"\nminCount = 5\n"), 0644)
		cfg, _ := config.Load(dir)
		t.Equal(cfg.Progress.MinCount, 5)
		t.End()
	})
}

func TestLoadMalformed(t *testing.T) {
	Test(t, "config: Load returns error for malformed toml", func(t *T) {
		dir := t.TB().TempDir()
		os.WriteFile(filepath.Join(dir, ".indra.toml"), []byte(`[bad`), 0644)
		_, err := config.Load(dir)
		t.Ok(err)

		t.End()
	})
}

func TestIsIgnoredExact(t *testing.T) {
	Test(t, "isIgnored: matches exact filename", func(t *T) {
		t.Ok(config.IsIgnored([]string{"foo.go"}, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredSingleStar(t *testing.T) {
	Test(t, "isIgnored: matches single star glob", func(t *T) {
		t.Ok(config.IsIgnored([]string{"*.go"}, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredDoubleStarPrefix(t *testing.T) {
	Test(t, "isIgnored: matches double star prefix", func(t *T) {
		t.Ok(config.IsIgnored([]string{"vendor/**"}, "vendor/pkg/file.go"))
		t.End()
	})
}

func TestIsIgnoredDoubleStarSuffix(t *testing.T) {
	Test(t, "isIgnored: matches double star suffix", func(t *T) {
		t.Ok(config.IsIgnored([]string{"**/*_test.go"}, "internal/foo/foo_test.go"))
		t.End()
	})
}

func TestIsIgnoredNoMatch(t *testing.T) {
	Test(t, "isIgnored: no match returns false", func(t *T) {
		t.Ok(!config.IsIgnored([]string{"vendor/**"}, "internal/foo.go"))
		t.End()
	})
}

func TestIsIgnoredEmpty(t *testing.T) {
	Test(t, "isIgnored: empty patterns returns false", func(t *T) {
		t.Ok(!config.IsIgnored(nil, "foo.go"))
		t.End()
	})
}

func TestIsIgnoredZeroSegments(t *testing.T) {
	Test(t, "isIgnored: double star matches zero segments", func(t *T) {
		t.Ok(config.IsIgnored([]string{"vendor/**"}, "vendor/file.go"))
		t.End()
	})
}

func TestIsIgnoredNegationOverrides(t *testing.T) {
	Test(t, "isIgnored: negation overrides earlier match", func(t *T) {
		t.Ok(!config.IsIgnored([]string{"vendor/**", "!vendor/mypkg/**"}, "vendor/mypkg/file.go"))
		t.End()
	})
}

func TestIsIgnoredPositiveAfterNegationWins(t *testing.T) {
	Test(t, "isIgnored: positive after negation wins", func(t *T) {
		t.Ok(config.IsIgnored([]string{"!vendor/**", "vendor/**"}, "vendor/pkg/file.go"))
		t.End()
	})
}

func TestIsIgnoredNegationOnlyMatch(t *testing.T) {
	Test(t, "isIgnored: negation-only match unignores", func(t *T) {
		t.Ok(!config.IsIgnored([]string{"!vendor/**"}, "vendor/pkg/file.go"))
		t.End()
	})
}

func TestDefaultIgnoreIgnoresVendor(t *testing.T) {
	Test(t, "default: ignores vendor", func(t *T) {
		t.Ok(config.IsIgnored(config.DefaultIgnorePatterns, "vendor/pkg/f.go"))
		t.End()
	})
}

func TestDefaultIgnoreIgnoresTestdata(t *testing.T) {
	Test(t, "default: ignores testdata", func(t *T) {
		t.Ok(config.IsIgnored(config.DefaultIgnorePatterns, "testdata/f.go"))
		t.End()
	})
}

func TestDefaultIgnoreIgnoresDotDir(t *testing.T) {
	Test(t, "default: ignores dot-dir", func(t *T) {
		t.Ok(config.IsIgnored(config.DefaultIgnorePatterns, ".git/config"))
		t.End()
	})
}

func TestDefaultIgnoreKeepsInternal(t *testing.T) {
	Test(t, "default: does not ignore internal", func(t *T) {
		t.Ok(!config.IsIgnored(config.DefaultIgnorePatterns, "internal/pkg/f.go"))
		t.End()
	})
}

func TestToLoaderConfigOnEnabled(t *testing.T) {
	Test(t, "ToLoaderConfig: on maps to Enabled true", func(t *T) {
		cfg := config.Config{Rules: map[string]string{"tape": "on"}}
		lc := cfg.ToLoaderConfig()
		t.Equal(lc["tape"].Enabled, true)
		t.End()
	})
}

func TestToLoaderConfigOffDisabled(t *testing.T) {
	Test(t, "ToLoaderConfig: off maps to Enabled false", func(t *T) {
		cfg := config.Config{Rules: map[string]string{"tape": "off"}}
		lc := cfg.ToLoaderConfig()
		t.Equal(lc["tape"].Enabled, false)
		t.End()
	})
}

func TestIsIgnoredDoubleStarExhausts(t *testing.T) {
	Test(t, "isIgnored: double star backtracking exhausts and returns false", func(t *T) {
		t.Ok(!config.IsIgnored([]string{"**/nomatch"}, "a/b/c"))
		t.End()
	})
}
