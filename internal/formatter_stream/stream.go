// Package formatter_stream renders findings like dump but with box-drawing
// column separators, producing a table with borders.
package formatter_stream

import (
	"fmt"
	"strings"
	"text/tabwriter"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

// Format is the stream formatter — same as dump but with │ column separators.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	last := index == count-1

	if len(places) == 0 {
		if !last || errorsCount == 0 {
			return ""
		}
		return dump.Summary(errorsCount, filesWithIssues)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "\033[4m%s\033[0m\n", name)

	tw := tabwriter.NewWriter(&sb, 0, 0, 1, ' ', 0)
	for _, pl := range places {
		fmt.Fprintf(tw, "  \033[90m%d:%d\033[0m\t │ \033[31merror\033[0m   %s\t │ \033[90m%s\033[0m\n",
			pl.Position.Line, pl.Position.Column, pl.Message, pl.Rule)
	}
	tw.Flush()

	if last {
		sb.WriteString(dump.Summary(errorsCount, filesWithIssues))
	}
	return sb.String()
}
