package lint

import (
	"fmt"
	"io"
	"os"

	"coderaiser/indra/internal/engine"
	"coderaiser/indra/internal/plugins"
)

func Run(files []string, w io.Writer) bool {
	return run(files, w, false, os.WriteFile)
}

func Fix(files []string, w io.Writer) bool {
	return run(files, w, true, os.WriteFile)
}

func run(files []string, w io.Writer, fix bool, writeFile func(string, []byte, os.FileMode) error) bool {
	failed := false

	for _, filename := range files {
		src, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(w, "file://%s: %v\n", filename, err)
			failed = true
			continue
		}

		out, places, err := engine.Indra(src, plugins.All, fix)
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

		if fix && string(out) != string(src) {
			if err := writeFile(filename, out, 0644); err != nil {
				fmt.Fprintf(w, "file://%s: fix: %v\n", filename, err)
			}
		}
	}

	return failed
}
