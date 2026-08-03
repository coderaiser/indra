package indra

import (
	"errors"
	"fmt"
	"io"
	"strings"

	loader "coderaiser/indra/engine-loader"
	processor "coderaiser/indra/engine-processor"
	runner "coderaiser/indra/engine-runner"
	"coderaiser/indra/internal/plugins"
	processor_go "coderaiser/indra/processor-go"
	"coderaiser/indra/types"
)

// loadPlugins resolves all enabled plugins into runnable items.
func loadPlugins() []runner.PluginItem {
	kinds := loader.Load(plugins.All, loader.DefaultConfig())
	items := make([]runner.PluginItem, len(kinds))
	for i, k := range kinds {
		items[i] = runner.PluginItem{Rule: k.Name(), Plugin: k}
	}
	return items
}

// Lint runs all plugins against src.
// Returns rewritten source, findings, and any parse error.
func Lint(src []byte, fix bool) ([]byte, []types.Place, error) {
	res, err := processor.Process(processor.Params{
		Src:     src,
		Fix:     fix,
		Plugins: loadPlugins(),
	})
	return res.Out, res.Places, err
}

// Indra runs the CLI over the given files, printing findings.
func Indra(args []string, w io.Writer) error {
	for _, arg := range args {
		if arg == "--version" || arg == "-v" {
			fmt.Fprintln(w, VersionLine())
			return nil
		}
	}

	fix := false
	files := []string{}
	for _, a := range args {
		if a == "--fix" {
			fix = true
			continue
		}
		if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}

	if len(files) == 0 {
		return nil
	}

	items := loadPlugins()
	failed := false
	for _, filename := range files {
		places, err := processor_go.ProcessFile(filename, processor_go.Opt(items, fix))
		if err != nil {
			fmt.Fprintf(w, "file://%s: %v\n", filename, err)
			failed = true
			continue
		}
		for _, pl := range places {
			failed = true
			fmt.Fprintf(w, "file://%s:%d:%d: %s\n",
				filename,
				pl.Pos.Line,
				pl.Pos.Column,
				pl.Message,
			)
		}
	}

	if failed {
		return errors.New("lint failed")
	}
	return nil
}
