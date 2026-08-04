package formatter_dump_test

import (
	"strings"
	"testing"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

var places1 = []types.Place{
	{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 10, Column: 3}},
}

func TestDumpNoPlacesMidRun(t *testing.T) {
	tape.Test(t, "dump: no places mid-run returns empty", func(t *tape.T) {
		out := dump.Format("foo.go", nil, 0, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestDumpNoPlacesLastFileNoErrors(t *testing.T) {
	tape.Test(t, "dump: no places last file no errors returns empty", func(t *tape.T) {
		out := dump.Format("foo.go", nil, 4, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestDumpNoPlacesLastFileWithErrors(t *testing.T) {
	tape.Test(t, "dump: no places last file with prior errors shows summary", func(t *tape.T) {
		out := dump.Format("foo.go", nil, 4, 5, 2, 3)
		t.Ok(strings.Contains(out, "3 errors"))
		t.End()
	})
}

func TestDumpWithPlacesShowsFilename(t *testing.T) {
	tape.Test(t, "dump: places shows filename", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 0, 5, 1, 1)
		t.Ok(strings.Contains(out, "foo.go"))
		t.End()
	})
}

func TestDumpWithPlacesShowsLineCol(t *testing.T) {
	tape.Test(t, "dump: places shows line:col", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 0, 5, 1, 1)
		t.Ok(strings.Contains(out, "10:3"))
		t.End()
	})
}

func TestDumpWithPlacesShowsMessage(t *testing.T) {
	tape.Test(t, "dump: places shows message", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 0, 5, 1, 1)
		t.Ok(strings.Contains(out, "remove Test.Skip call"))
		t.End()
	})
}

func TestDumpWithPlacesShowsRule(t *testing.T) {
	tape.Test(t, "dump: places shows rule", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 0, 5, 1, 1)
		t.Ok(strings.Contains(out, "tape/remove-skip"))
		t.End()
	})
}

func TestDumpSummarySingular(t *testing.T) {
	tape.Test(t, "dump: summary singular error and file", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})
}

func TestDumpSummaryPlural(t *testing.T) {
	tape.Test(t, "dump: summary plural errors and files", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 4, 5, 2, 3)
		t.Ok(strings.Contains(out, "3 errors in 2 files"))
		t.End()
	})
}

func TestDumpSummaryFixHint(t *testing.T) {
	tape.Test(t, "dump: summary includes fix hint", func(t *tape.T) {
		out := dump.Format("foo.go", places1, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "--fix"))
		t.End()
	})
}
