package formatter_memory_test

import (
	"strings"
	"testing"

	formmem "coderaiser/indra/internal/formatter_memory"
	"coderaiser/indra/types"
	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}

func TestMemory(t *testing.T) {
	Test(t, "memory: mid-run returns empty", func(t *T) {
		out := formmem.Format("a.go", nil, nil, 0, 3, 0, 0)
		t.Equal(out, "")
		t.End()
	})

	Test(t, "memory: last file with issues contains dump summary", func(t *T) {
		out := formmem.Format("a.go", nil, []types.Place{place1}, 2, 3, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "memory: last file contains Memory section", func(t *T) {
		out := formmem.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "Memory"))
		t.End()
	})

	Test(t, "memory: last file contains HeapAlloc", func(t *T) {
		out := formmem.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "HeapAlloc"))
		t.End()
	})

	Test(t, "memory: last file contains Sys", func(t *T) {
		out := formmem.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "Sys"))
		t.End()
	})

	Test(t, "memory: last file contains NumGC", func(t *T) {
		out := formmem.Format("a.go", nil, nil, 2, 3, 0, 0)
		t.Ok(strings.Contains(out, "NumGC"))
		t.End()
	})
}
