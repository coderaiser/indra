//go:build linux || darwin

package formatter_progress_bar

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestTermWidthIoctlCol(t *testing.T) {
	Test(t, "term-unix: TermWidth uses ioctl column size", func(t *T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "")
		old := winsizeReader
		winsizeReader = func() (row, col uint16) { return 24, 160 }
		defer func() { winsizeReader = old }()
		result := TermWidth()
		t.Equal(result, 160)

		t.End()
	})

	Test(t, "term-unix: TermWidth falls back to 80 on zero column", func(t *T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "")
		old := winsizeReader
		winsizeReader = func() (row, col uint16) { return 24, 0 }
		defer func() { winsizeReader = old }()
		result := TermWidth()
		t.Equal(result, 80)

		t.End()
	})

	Test(t, "term-unix: ioctlWinsize returns row and col", func(t *T) {
		row, col := ioctlWinsize()
		t.Equal(row, col)
		t.End()
	})
}
