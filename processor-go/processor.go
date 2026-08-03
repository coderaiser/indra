// Package processor_go processes Go source files and directories.
package processor_go

import (
	"os"
	"path/filepath"
	"strings"

	processor "coderaiser/indra/engine-processor"
	runner "coderaiser/indra/engine-runner"
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

// ProcessDir runs ProcessFile on every .go file in dir (non-recursive).
func ProcessDir(dir string, opts Options, ov ...Overrides) ([]types.Place, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var all []types.Place
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		places, err := ProcessFile(filepath.Join(dir, e.Name()), opts, ov...)
		if err != nil {
			return all, err
		}
		all = append(all, places...)
	}
	return all, nil
}
