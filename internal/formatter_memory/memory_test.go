package formatter_memory_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	formatter_memory "coderaiser/indra/internal/formatter_memory"
	pb "coderaiser/indra/internal/formatter_progress_bar"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}

func TestMemory(t *testing.T) {
	Test(t, "memory: mid-run returns empty", func(t *T) {
		out := formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "memory: last file with issues contains dump summary", func(t *T) {
		out := formatter_memory.Format("a.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "memory: last file contains Memory section", func(t *T) {
		out := formatter_memory.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "Memory"))
		t.End()
	})

	Test(t, "memory: last file contains HeapAlloc", func(t *T) {
		out := formatter_memory.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "HeapAlloc"))
		t.End()
	})

	Test(t, "memory: last file contains Sys", func(t *T) {
		out := formatter_memory.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "Sys"))
		t.End()
	})

	Test(t, "memory: last file contains NumGC", func(t *T) {
		out := formatter_memory.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "NumGC"))
		t.End()
	})
}

func TestMemoryFallback(t *testing.T) {
	Test(t, "memory fallback: mid-run returns empty when bar hidden", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "memory fallback: last file contains Memory when bar hidden", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := formatter_memory.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "Memory"))
		t.End()
	})
}

func TestMemoryProgressBar(t *testing.T) {
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

	Test(t, "memory progress bar: mid-run writes hide cursor to stderr", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_memory.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, pb.HideCursor))
		t.End()
	})

	Test(t, "memory progress bar: mid-run writes show cursor to stderr", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_memory.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, pb.ShowCursor))
		t.End()
	})

	Test(t, "memory progress bar: mid-run bar contains heap stats", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		out := captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
			formatter_memory.Format("c.go", nil, nil, 2, 3, 0, 0)
		})
		t.Ok(strings.Contains(out, "heap"))
		t.End()
	})

	Test(t, "memory progress bar: last file returns dump summary", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
		})
		out := formatter_memory.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "memory progress bar: last file contains Memory block", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
		})
		out := formatter_memory.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "Memory"))
		t.End()
	})

	Test(t, "memory progress bar: last file contains HeapAlloc", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")

		captureStderr(t, func() {
			formatter_memory.Format("a.go", nil, nil, 0, 3, 0, 0)
			formatter_memory.Format("b.go", nil, nil, 1, 3, 0, 0)
		})
		out := formatter_memory.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "HeapAlloc"))
		t.End()
	})
}
