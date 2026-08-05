// Package engine_processor orchestrates the parse → run → print pipeline.
package engine_processor

import (
	engineparser "coderaiser/indra/engine_parser"
	engineprinter "coderaiser/indra/engine_printer"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/types"
)

// Params are the inputs to Process.
type Params struct {
	Src      []byte
	Filename string
	Fix      bool
	Plugins  []runner.PluginItem
}

// Result holds the processed output.
type Result struct {
	Out    []byte
	Places []types.Place
}

// print formats an AST file back to source bytes. It is a variable so tests
// can exercise the printer error path.
var print = engineprinter.Print

// Process parses Src, runs plugins, and prints the result if fix is set.
func Process(p Params) (Result, error) {
	file, fset, err := engineparser.Parse(p.Src)
	if err != nil {
		return Result{Out: p.Src}, err
	}
	places := runner.RunPlugins(runner.RunParams{
		File:    file,
		Fset:    fset,
		Fix:     p.Fix,
		Plugins: p.Plugins,
	})
	if p.Fix {
		out, err := print(file, fset)
		if err != nil {
			return Result{Out: p.Src, Places: places}, err
		}
		return Result{Out: out, Places: places}, nil
	}
	return Result{Out: p.Src, Places: places}, nil
}
