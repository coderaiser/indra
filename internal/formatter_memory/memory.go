// Package formatter_memory wraps dump and appends heap statistics on the
// final file, with an optional live progress bar.
package formatter_memory

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	pb "coderaiser/indra/internal/formatter_progress_bar"
	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

// Format wraps dump and appends memory stats on the last file.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	result := dump.Format(name, nil, places, index, count, filesWithIssues, errorsCount)

	if !pb.ShouldShow(count) {
		if index != count-1 {
			return ""
		}
		return result + memoryBlock()
	}

	if index == 0 {
		fmt.Fprintf(os.Stderr, "%s", pb.HideCursor)
	}

	errStr := "👌"
	if errorsCount > 0 {
		errStr = fmt.Sprintf("\033[31m%d\033[0m", errorsCount)
	}

	bar := pb.RenderBar(index+1, count, "#ea4335")
	pct := 0
	if count > 0 {
		pct = (index + 1) * 100 / count
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	heapLine := fmt.Sprintf("heap %s / %s  sys %s",
		formatBytes(ms.HeapAlloc), formatBytes(ms.HeapSys), formatBytes(ms.Sys))

	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s | %s",
		bar, pct, errStr, index+1, count, pb.Truncate(name, 30), heapLine)
	width := pb.TermWidth()
	if pb.VisibleLen(line) > width {
		line = pb.TruncateANSI(line, width)
	}
	fmt.Fprintf(os.Stderr, "\r%s", line)

	if index == count-1 {
		fmt.Fprintf(os.Stderr, "\r\033[2K%s", pb.ShowCursor)
		return result + "\n" + memoryBlock()
	}
	return ""
}

func memoryBlock() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	var sb strings.Builder
	sb.WriteString("\nMemory:\n")
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
