package formatter_time_test

import (
	"strings"
	"testing"

	formtime "coderaiser/indra/internal/formatter_time"
	"coderaiser/indra/types"
	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}

func TestTime(t *testing.T) {
	Test(t, "time: mid-run returns empty", func(t *T) {
		out := formtime.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.Equal(out, "")
		t.End()
	})

	Test(t, "time: last file with issues contains dump summary", func(t *T) {
		formtime.Format("a.go", nil, nil, 0, 3, 0, 0)
		formtime.Format("b.go", nil, nil, 1, 3, 0, 0)
		out := formtime.Format("c.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "time: last file contains Time section", func(t *T) {
		out := formtime.Format("a.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.Contains(out, "Time:"))
		t.End()
	})

	Test(t, "time: elapsed is non-negative", func(t *T) {
		formtime.Format("a.go", nil, nil, 0, 2, 0, 0)
		out := formtime.Format("b.go", nil, nil, 1, 2, 0, 0)
		t.Ok(strings.Contains(out, "ms") || strings.Contains(out, "µs") || strings.Contains(out, "s"))
		t.End()
	})

	Test(t, "time: state resets between runs", func(t *T) {
		formtime.Format("a.go", nil, nil, 0, 1, 0, 0)
		out := formtime.Format("a.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.Contains(out, "Time:"))
		t.End()
	})
}
