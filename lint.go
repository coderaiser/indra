package indra

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"coderaiser/indra/internal/lint"
)

func filterFlags(args []string) []string {
	var files []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	return files
}

func RunLint(args []string, w io.Writer) error {
	for _, arg := range args {
		if arg == "--version" || arg == "-v" {
			fmt.Fprintln(w, VersionLine())
			return nil
		}
		if arg == "--help" || arg == "-h" {
			fmt.Fprint(w, Help())
			return nil
		}
	}

	files := filterFlags(args)

	if len(files) == 0 {
		return nil
	}

	failed := lint.Run(files, w)

	if failed {
		return errors.New("lint failed")
	}

	return nil
}