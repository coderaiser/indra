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
// Tape sub-rules fire only under their "tape/*" group names; the standalone
// provider entries are stripped so each rule fires exactly once.
func loadPlugins() []runner.PluginItem {
	kinds := loader.Load(plugins.LoadInput(), loader.DefaultConfig())
	items := make([]runner.PluginItem, 0, len(kinds))
	for _, k := range kinds {
		if isProviderName(k.Name()) {
			continue
		}
		items = append(items, runner.PluginItem{Rule: k.Name(), Plugin: k})
	}
	return items
}

// isProviderName reports whether name is a standalone tape sub-rule that should
// fire only via its tape group.
func isProviderName(name string) bool {
	for _, p := range plugins.Providers {
		if p.Name == name {
			return true
		}
	}
	return false
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
				pl.Position.Line,
				pl.Position.Column,
				pl.Message,
			)
		}
	}

	if failed {
		return errors.New("lint failed")
	}
	return nil
}
