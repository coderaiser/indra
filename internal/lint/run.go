package lint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"

	"coderaiser/indra/internal/lint/rule"
	"coderaiser/indra/internal/lint/rules"
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
		fset := token.NewFileSet()

		file, err := parser.ParseFile(
			fset,
			filename,
			nil,
			parser.ParseComments,
		)

		if err != nil {
			fmt.Fprintf(
				w,
				"file://%s: %v\n",
				filename,
				err,
			)

			failed = true
			continue
		}

		modified := false

		for _, r := range rules.All {
			if fix {
				if fixer, ok := r.(rule.Fixer); ok {
					if fixer.Fix(file, fset) {
						modified = true
					}
				}
			}

			results := r.Check(file, fset)

			for _, result := range results {
				failed = true

				fmt.Fprintf(
					w,
					"file://%s:%d:%d: %s\n",
					result.Pos.Filename,
					result.Pos.Line,
					result.Pos.Column,
					result.Message,
				)
			}
		}

		if modified {
			if err := writeFormatted(filename, file, fset, writeFile); err != nil {
				fmt.Fprintf(w, "file://%s: fix: %v\n", filename, err)
			}
		}
	}

	return failed
}

func writeFormatted(filename string, file *ast.File, fset *token.FileSet, writeFile func(string, []byte, os.FileMode) error) error {
	var buf bytes.Buffer
	format.Node(&buf, fset, file)
	return writeFile(filename, buf.Bytes(), 0644)
}

