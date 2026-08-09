package formatter_stream_test

import (
	"strings"
	"testing"

	formatter_stream "coderaiser/indra/internal/formatter_stream"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 5, Column: 2}}

func TestStream(t *testing.T) {
	Test(t, "stream: no places mid-run returns empty", func(t *T) {
		out := formatter_stream.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "stream: places shows filename", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 3, 1, 1)
		t.Ok(strings.Contains(out, "a.go"))
		t.End()
	})

	Test(t, "stream: places shows line:col", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 3, 1, 1)
		t.Ok(strings.Contains(out, "5:2"))
		t.End()
	})

	Test(t, "stream: places shows message", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 3, 1, 1)
		t.Ok(strings.Contains(out, "remove Test.Skip call"))
		t.End()
	})

	Test(t, "stream: places shows rule", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 3, 1, 1)
		t.Ok(strings.Contains(out, "tape/remove-skip"))
		t.End()
	})

	Test(t, "stream: last file with errors appends summary", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "stream: last file clean no errors returns empty", func(t *T) {
		out := formatter_stream.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "stream: last file clean with prior errors returns summary only", func(t *T) {
		out := formatter_stream.Format("a.go", nil, nil, 2, 3, 1, 2)
		t.Ok(strings.Contains(out, "2 errors"))
		t.End()
	})

	Test(t, "stream: first of two files with errors shows filename", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		t.Ok(strings.Contains(out, "a.go"))
		t.End()
	})

	Test(t, "stream: first of two files with errors shows line:col", func(t *T) {
		out := formatter_stream.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		t.Ok(strings.Contains(out, "5:2"))
		t.End()
	})
}
