// Package processor_go processes Go source files and directories.
package processor_go

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	processor "coderaiser/indra/engine-processor"
	runner "coderaiser/indra/engine-runner"
	"coderaiser/indra/internal/config"
	"coderaiser/indra/types"
)

// Options configures a ProcessFile or ProcessDir call.
type Options struct {
	plugins []runner.PluginItem
	fix     bool
}

// Opt returns an Options with the given plugins and fix flag.
func Opt(plugins []runner.PluginItem, fix bool) Options {
	return Options{plugins: plugins, fix: fix}
}

// Overrides replaces internal functions — for testing only.
type Overrides struct {
	writeFile func(string, []byte, os.FileMode) error
}

// WithWriteFile returns an Overrides with a custom writeFile func.
func WithWriteFile(fn func(string, []byte, os.FileMode) error) Overrides {
	return Overrides{writeFile: fn}
}

// ProcessFile reads path, runs plugins, and writes back if fix=true and changed.
// Skips files containing "//go:build ignore".
func ProcessFile(path string, opts Options, ov ...Overrides) ([]types.Place, error) {
	writeFn := os.WriteFile
	if len(ov) > 0 && ov[0].writeFile != nil {
		writeFn = ov[0].writeFile
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(src), "//go:build ignore") {
		return nil, nil
	}
	res, err := processor.Process(processor.Params{
		Src:      src,
		Filename: path,
		Fix:      opts.fix,
		Plugins:  opts.plugins,
	})
	if err != nil {
		return nil, err
	}
	if opts.fix && string(res.Out) != string(src) {
		if err := writeFn(path, res.Out, 0644); err != nil {
			return nil, err
		}
	}
	return res.Places, nil
}

// ProcessDir walks dir recursively, running ProcessFile on every .go file.
// Skips: vendor/, testdata/, hidden dirs (dot-prefix), build-ignored files, and
// paths matching the ignore glob patterns (relative to dir).
// ignore may be nil to skip no ignore matching.
func ProcessDir(dir string, opts Options, ignore []string, ov ...Overrides) ([]types.Place, error) {
	var all []types.Place
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if path == dir {
				return nil
			}
			if name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		if config.IsIgnored(ignore, rel) {
			return nil
		}
		places, err := ProcessFile(path, opts, ov...)
		if err != nil {
			return err
		}
		all = append(all, places...)
		return nil
	})
	if err != nil {
		return all, err
	}
	return all, nil
}

// ResolveArgs splits CLI args into individual files and directories.
// Handles file paths, dir paths, and ./... / ... suffixes.
func ResolveArgs(args []string) (files []string, dirs []string) {
	for _, a := range args {
		if strings.HasSuffix(a, "/...") || a == "..." {
			dir := strings.TrimSuffix(a, "/...")
			if dir == "" {
				dir = "."
			}
			dirs = append(dirs, dir)
			continue
		}
		if strings.HasSuffix(a, "...") {
			dirs = append(dirs, strings.TrimSuffix(a, "..."))
			continue
		}
		info, err := os.Stat(a)
		if err != nil || !info.IsDir() {
			files = append(files, a)
			continue
		}
		dirs = append(dirs, a)
	}
	return files, dirs
}

// CollectFiles returns all .go files ResolveArgs would process: the explicit
// files plus every discovered file under directories.
func CollectFiles(files []string, dirs []string, ignore []string) []string {
	result := append([]string{}, files...)
	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			name := d.Name()
			if d.IsDir() {
				if path == dir {
					return nil
				}
				if name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(name, ".go") {
				return nil
			}
			rel, _ := filepath.Rel(dir, path)
			if config.IsIgnored(ignore, rel) {
				return nil
			}
			result = append(result, path)
			return nil
		})
	}
	return result
}
