// Package formatter_time wraps dump and appends elapsed duration on the
// final file, with an optional live progress bar.
package formatter_time

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	dump "coderaiser/indra/internal/formatter_dump"
	pb "coderaiser/indra/internal/formatter_progress_bar"
	"coderaiser/indra/types"
)

var (
	mu    sync.Mutex
	start time.Time
)

// Format wraps dump and appends elapsed time on the last file.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	mu.Lock()
	defer mu.Unlock()

	if index == 0 {
		start = time.Now()
	}

	result := dump.Format(name, nil, places, index, count, filesWithIssues, errorsCount)

	if !pb.ShouldShow(count) {
		if index != count-1 {
			return ""
		}
		return result + formatTimeBlock(time.Since(start))
	}

	if index == 0 {
		fmt.Fprintf(os.Stderr, "%s", pb.HideCursor)
	}

	errStr := "👌"
	if errorsCount > 0 {
		errStr = fmt.Sprintf("\033[31m%d\033[0m", errorsCount)
	}

	bar := pb.RenderBar(index+1, count, "#91e0d5")
	pct := 0
	if count > 0 {
		pct = (index + 1) * 100 / count
	}

	elapsed := time.Since(start)
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s | ⏳ %s",
		bar, pct, errStr, index+1, count, pb.Truncate(name, 30), formatDuration(elapsed))
	width := pb.TermWidth()
	if pb.VisibleLen(line) > width {
		line = pb.TruncateANSI(line, width)
	}
	fmt.Fprintf(os.Stderr, "\r%s", line)

	if index == count-1 {
		elapsed := time.Since(start)
		fmt.Fprintf(os.Stderr, "\r\033[2K%s", pb.ShowCursor)
		start = time.Time{} // reset
		return result + "\n" + formatTimeBlock(elapsed)
	}
	return ""
}

func formatTimeBlock(elapsed time.Duration) string {
	var sb strings.Builder
	sb.WriteString("\nTime: ")
	sb.WriteString(formatDuration(elapsed))
	sb.WriteString("\n")
	return sb.String()
}

func formatDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%dms", d.Milliseconds())
	default:
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
}
