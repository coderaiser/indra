package formatter_progress_test

import (
	"strings"
	"testing"

	formatter_progress "coderaiser/indra/internal/formatter_progress"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}

func TestProgress(t *testing.T) {
	Test(t, "progress: mid-run returns empty when below minCount", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "10")
		out := formatter_progress.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "progress: last file below minCount returns dump output", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "10")
		out := formatter_progress.Format("a.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error"))
		t.End()
	})

	Test(t, "progress: mid-run at or above minCount returns empty", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "0")
		out := formatter_progress.Format("a.go", nil, nil, 0, 5, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "progress: last file at or above minCount returns dump", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "0")
		out := formatter_progress.Format("a.go", nil, []types.Place{place1}, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error"))
		t.End()
	})
}

func TestProgressShouldShow(t *testing.T) {
	Test(t, "progress: ShouldShow uses default min when env unset", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "")
		t.Ok(formatter_progress.ShouldShow(1))
		t.End()
	})

	Test(t, "progress: ShouldShow false when count not above default min", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "")
		t.NotOk(formatter_progress.ShouldShow(0))

		t.End()
	})

	Test(t, "progress: ShouldShow honors valid min env", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "5")
		t.Ok(formatter_progress.ShouldShow(6))
		t.End()
	})

	Test(t, "progress: ShouldShow hides count not above valid min env", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "5")
		t.NotOk(formatter_progress.ShouldShow(5))

		t.End()
	})

	Test(t, "progress: ShouldShow ignores invalid min env", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_MIN", "abc")
		t.Ok(formatter_progress.ShouldShow(1))
		t.End()
	})
}
