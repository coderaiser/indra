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

// ProcessFile reads path, runs plugins, and writes back if fix=true and changed.
// Skips files containing "//go:build ignore".
func ProcessFile(path string, plugins []runner.PluginItem, fix bool) ([]types.Place, error) {
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
		Fix:      fix,
		Plugins:  plugins,
	})
	if err != nil {
		return nil, err
	}
	if fix && string(res.Out) != string(src) {
		if err := os.WriteFile(path, res.Out, 0644); err != nil {
			return nil, err
		}
	}
	return res.Places, nil
}

// ProcessDir runs ProcessFile on every .go file in dir (non-recursive).
func ProcessDir(dir string, plugins []runner.PluginItem, fix bool) ([]types.Place, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var all []types.Place
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		places, err := ProcessFile(filepath.Join(dir, e.Name()), plugins, fix)
		if err != nil {
			return all, err
		}
		all = append(all, places...)
	}
	return all, nil
}
