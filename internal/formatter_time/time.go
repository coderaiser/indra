// Package formatter_time wraps dump and appends elapsed duration on the
// final file.
package formatter_time

import (
	"fmt"
	"strings"
	"sync"
	"time"

	dump "coderaiser/indra/internal/formatter_dump"
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

	if index != count-1 {
		return ""
	}

	elapsed := time.Since(start)
	start = time.Time{} // reset

	var sb strings.Builder
	sb.WriteString(result)
	fmt.Fprintf(&sb, "\nTime: %s\n", formatDuration(elapsed))
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
