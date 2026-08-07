// Package formatter_frame wraps the codeframe formatter with a percent
// progress indicator written to stderr mid-run.
package formatter_frame

import (
	"fmt"
	"math"
	"os"

	codeframe "coderaiser/indra/internal/formatter_codeframe"
	"coderaiser/indra/types"
)

// Format wraps codeframe with \r42% progress on stderr mid-run.
func Format(name string, source []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	result := codeframe.Format(name, source, places, index, count, filesWithIssues, errorsCount)
	last := index == count-1

	pct := int(math.Round(float64(index+1) / float64(count) * 100))
	fmt.Fprintf(os.Stderr, "\r%d%%", pct)

	if last {
		fmt.Fprintf(os.Stderr, "\r\033[2K")
		return result
	}
	return ""
}
