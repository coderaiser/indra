// Package formatter_dump renders lint findings as a plain-text dump.
package formatter_dump

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"coderaiser/indra/types"
)

// Format is the dump formatter — called once per file with running totals.
// It returns the rendered block for the given file, or "" when there is
// nothing to print (no places mid-run, or a clean last file). On the last
// file it appends a summary of errors.
func Format(name string, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	last := index == count-1

	if len(places) == 0 {
		if !last || errorsCount == 0 {
			return ""
		}
		return summary(errorsCount, filesWithIssues)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "\033[4m%s\033[0m\n", name)

	tw := tabwriter.NewWriter(&sb, 0, 0, 3, ' ', 0)
	for _, pl := range places {
		fmt.Fprintf(tw, "  \033[90m%d:%d\033[0m\t\033[31merror\033[0m   %s\t\033[90m%s\033[0m\n",
			pl.Position.Line, pl.Position.Column, pl.Message, pl.Rule)
	}
	tw.Flush()

	if last {
		sb.WriteString(summary(errorsCount, filesWithIssues))
	}
	return sb.String()
}

func summary(errorsCount, filesWithIssues int) string {
	errWord := "errors"
	if errorsCount == 1 {
		errWord = "error"
	}
	fileWord := "files"
	if filesWithIssues == 1 {
		fileWord = "file"
	}
	return fmt.Sprintf("\033[1;91m✖ %d %s in %d %s\033[0m\n  fixable with the --fix option\n",
		errorsCount, errWord, filesWithIssues, fileWord)
}
