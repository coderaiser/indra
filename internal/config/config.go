// Package config parses the .indra.toml configuration file and provides
// ignore-pattern matching.
package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the parsed .indra.toml.
type Config struct {
	Rules    map[string]string `toml:"rules"`
	Ignore   IgnoreConfig      `toml:"ignore"`
	Progress ProgressConfig    `toml:"progress"`
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

// IsIgnored reports whether relPath matches any ignore pattern.
// relPath is slash-separated and relative to the walk root.
func IsIgnored(patterns []string, relPath string) bool {
	relPath = filepath.ToSlash(relPath)
	for _, pat := range patterns {
		if matchGlob(pat, relPath) {
			return true
		}
	}
	return false
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
