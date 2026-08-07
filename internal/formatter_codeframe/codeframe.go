// Package formatter_codeframe shows each error with surrounding source
// context, highlighting the offending line.
package formatter_codeframe

import (
	"bytes"
	"fmt"
	"strings"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

const context = 2 // lines before and after

// Format shows each error with surrounding source context.
// Falls back to dump if source is nil or line is out of range.
func Format(name string, source []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	last := index == count-1

	if len(places) == 0 {
		if !last || errorsCount == 0 {
			return ""
		}
		return dump.Summary(errorsCount, filesWithIssues)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "\033[4m%s\033[0m\n", name)

	lines := splitLines(source)

	for _, pl := range places {
		fmt.Fprintf(&sb, "  \033[90m%d:%d\033[0m  \033[31merror\033[0m   %s   \033[90m%s\033[0m\n",
			pl.Position.Line, pl.Position.Column, pl.Message, pl.Rule)

		col := pl.Position.Column
		if col < 1 {
			col = 1
		}
		indent := strings.Repeat(" ", col-1)
		numWidth := len(fmt.Sprintf("%d", pl.Position.Line))
		prefix := strings.Repeat(" ", numWidth+2+1)
		fmt.Fprintf(&sb, "%s| %s\033[33m^\033[0m %s \033[90m(%s)\033[0m\n",
			prefix, indent, pl.Message, pl.Rule)

		if source == nil || pl.Position.Line < 1 || pl.Position.Line > len(lines) {
			continue
		}

		start := pl.Position.Line - context - 1
		if start < 0 {
			start = 0
		}
		end := pl.Position.Line + context
		if end > len(lines) {
			end = len(lines)
		}

		for i := start; i < end; i++ {
			lineNum := i + 1
			marker := "  "
			color := "\033[0m"
			if lineNum == pl.Position.Line {
				marker = "> "
				color = "\033[1m"
			}
			fmt.Fprintf(&sb, "%s\033[90m%d\033[0m │ %s%s\033[0m\n",
				marker, lineNum, color, lines[i])
		}
		sb.WriteByte('\n')
	}

	if last {
		sb.WriteString(dump.Summary(errorsCount, filesWithIssues))
	}
	return sb.String()
}

func splitLines(src []byte) []string {
	if src == nil {
		return nil
	}
	lines := bytes.Split(src, []byte("\n"))
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = string(l)
	}
	return result
}
