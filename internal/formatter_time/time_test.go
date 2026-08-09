package formatter_time_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	pb "coderaiser/indra/internal/formatter_progress_bar"
	formatter_time "coderaiser/indra/internal/formatter_time"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}

func TestTime(t *testing.T) {
	Test(t, "time: mid-run returns empty", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "time: last file with issues contains dump summary", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
		formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
		out := formatter_time.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "time: last file contains Time section", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := formatter_time.Format("a.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.Contains(out, "Time:"))
		t.End()
	})

	Test(t, "time: elapsed is non-negative", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		formatter_time.Format("a.go", nil, nil, 0, 2, 0, 0)
		out := formatter_time.Format("b.go", nil, nil, 1, 2, 0, 0)
		t.Ok(strings.Contains(out, "ms") || strings.Contains(out, "µs") || strings.Contains(out, "s"))
		t.End()
	})

	Test(t, "time: state resets between runs", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		formatter_time.Format("a.go", nil, nil, 0, 1, 0, 0)
		out := formatter_time.Format("a.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.Contains(out, "Time:"))
		t.End()
	})
}

func TestTimeProgressBar(t *testing.T) {
	captureStderr := func(t *T, fn func()) string {
		t.TB().Helper()
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		t.TB().Cleanup(func() {
			w.Close()
			os.Stderr = oldStderr
			r.Close()
		})
		fn()
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	Test(t, "time progress bar: mid-run writes hide cursor to stderr", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_time.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, pb.HideCursor))
		t.End()
	})

	Test(t, "time progress bar: mid-run writes show cursor to stderr", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_time.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, pb.ShowCursor))
		t.End()
	})

	Test(t, "time progress bar: mid-run bar contains elapsed time", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_time.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, "⏳"))
		t.End()
	})

	Test(t, "time progress bar: last file returns dump summary", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		captureStderr(t, func() {
			formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
		})
		out := formatter_time.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "time progress bar: last file contains Time block", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		captureStderr(t, func() {
			formatter_time.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 3, 0, 0)
		})
		out := formatter_time.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "Time:"))
		t.End()
	})

	Test(t, "time progress bar: truncates overlong line to term width", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "10")

		captureStderr(t, func() {
			formatter_time.Format("a-very-long-filename-indeed.go", nil, nil, 0, 2, 0, 0)
			formatter_time.Format("b.go", nil, nil, 1, 2, 0, 0)
		})
		t.End()
	})
}
