// Package formatter_progress renders a simple percent progress on stderr
// mid-run and falls back to the dump format on the final file.
package formatter_progress

import (
	"fmt"
	"math"
	"os"
	"strconv"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

const defaultMin = 0

// Format is the progress formatter.
// Writes \r42% to stderr mid-run; returns dump output on last file.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	result := dump.Format(name, nil, places, index, count, filesWithIssues, errorsCount)
	last := index == count-1

	if !shouldShow(count) {
		return result
	}

	pct := int(math.Round(float64(index+1) / float64(count) * 100))
	fmt.Fprintf(os.Stderr, "\r%d%%", pct)

	if last {
		fmt.Fprintf(os.Stderr, "\r\033[2K")
		return result
	}
	return ""
}

// ShouldShow returns whether progress should be shown. Exported for testing.
func ShouldShow(count int) bool {
	min := defaultMin
	if v := os.Getenv("INDRA_PROGRESS_MIN"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			min = n
		}
	}
	return count > min
}

func shouldShow(count int) bool {
	return ShouldShow(count)
}
