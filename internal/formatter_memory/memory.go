// Package formatter_memory wraps dump and appends heap statistics on the
// final file.
package formatter_memory

import (
	"fmt"
	"runtime"
	"strings"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

// Format wraps dump and appends memory stats on the last file.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	result := dump.Format(name, nil, places, index, count, filesWithIssues, errorsCount)

	if index != count-1 {
		return ""
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	var sb strings.Builder
	sb.WriteString(result)
	fmt.Fprintf(&sb, "\nMemory:\n")
	fmt.Fprintf(&sb, "  HeapAlloc:   %s\n", formatBytes(ms.HeapAlloc))
	fmt.Fprintf(&sb, "  TotalAlloc:  %s\n", formatBytes(ms.TotalAlloc))
	fmt.Fprintf(&sb, "  Sys:         %s\n", formatBytes(ms.Sys))
	fmt.Fprintf(&sb, "  NumGC:       %d\n", ms.NumGC)
	return sb.String()
}

func formatBytes(b uint64) string {
	const mb = 1 << 20
	const kb = 1 << 10
	switch {
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/mb)
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/kb)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
