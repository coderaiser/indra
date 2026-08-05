// Package config parses the .indra.toml configuration file and provides
// ignore-pattern matching.
package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	loader "coderaiser/indra/engine_loader"

	"github.com/BurntSushi/toml"
)

// Config is the parsed .indra.toml.
type Config struct {
	Rules    map[string]string `toml:"rules"`
	Plugins  []string          `toml:"plugins"`
	Ignore   IgnoreConfig      `toml:"ignore"`
	Progress ProgressConfig    `toml:"progress"`
	Match    MatchConfig       `toml:"match"`
}

// MatchConfig maps a file glob pattern to per-rule on/off overrides. It backs
// the [match] section: for files matching the pattern, each listed rule is
// switched on or off regardless of the global [rules] value.
type MatchConfig map[string]map[string]string

// OverrideRules returns the merged rule overrides that apply to filename.
// Patterns are matched against the file's base name after sorting for
// deterministic order; when several patterns match the same rule, the last
// (alphabetically) matching pattern wins.
func (m MatchConfig) OverrideRules(filename string) map[string]string {
	out := map[string]string{}
	if len(m) == 0 {
		return out
	}
	base := filepath.Base(filename)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, pat := range keys {
		ok, err := filepath.Match(pat, base)
		if err != nil || !ok {
			continue
		}
		for rule, val := range m[pat] {
			out[rule] = val
		}
	}
	return out
}

// IgnoreConfig holds the ignore patterns.
type IgnoreConfig struct {
	Patterns []string `toml:"patterns"`
}

// ProgressConfig holds progress-bar display options from [progress] in .indra.toml.
type ProgressConfig struct {
	Color    string `toml:"color"`
	MinCount int    `toml:"minCount"`
}

// Load reads .indra.toml from dir. Returns a zero Config if the file is absent.
func Load(dir string) (Config, error) {
	var cfg Config
	path := filepath.Join(dir, ".indra.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}

// IsIgnored reports whether relPath is ignored by the pattern list.
// Patterns are evaluated in order; the last matching pattern wins.
// A pattern prefixed with "!" negates: the path is not ignored if the
// last matching pattern is a negation.
// relPath must be slash-separated and relative to the walk root.
func IsIgnored(patterns []string, relPath string) bool {
	relPath = filepath.ToSlash(relPath)
	ignored := false
	for _, pat := range patterns {
		negate := strings.HasPrefix(pat, "!")
		if negate {
			pat = pat[1:]
		}
		if matchGlob(pat, relPath) {
			ignored = !negate
		}
	}
	return ignored
}

// DefaultIgnorePatterns are the built-in base patterns, replacing the
// previously hardcoded vendor/testdata/dot-dir skips in processor-go.
var DefaultIgnorePatterns = []string{
	"vendor/**",
	"testdata/**",
	".*/**",
}

// ToLoaderConfig translates the [rules] section into a loader.Config.
// "on" → enabled, "off" → disabled.
func (c Config) ToLoaderConfig() loader.Config {
	lc := make(loader.Config, len(c.Rules))
	for rule, val := range c.Rules {
		lc[rule] = loader.RuleState{Enabled: val == "on"}
	}
	return lc
}

func matchGlob(pattern, path string) bool {
	return matchSegments(strings.Split(pattern, "/"), strings.Split(path, "/"))
}

func matchSegments(pat, path []string) bool {
	for len(pat) > 0 && len(path) > 0 {
		if pat[0] == "**" {
			pat = pat[1:]
			for i := 0; i <= len(path); i++ {
				if matchSegments(pat, path[i:]) {
					return true
				}
			}
			return false
		}
		ok, err := filepath.Match(pat[0], path[0])
		if err != nil || !ok {
			return false
		}
		pat = pat[1:]
		path = path[1:]
	}
	return len(pat) == 0 && len(path) == 0
}
