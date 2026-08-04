//go:build windows

package formatter_progress_bar

import (
	"os"
	"strconv"
)

// TermWidth returns the terminal width from INDRA_TERM_WIDTH env or 80.
func TermWidth() int {
	if v := os.Getenv("INDRA_TERM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 80
}
