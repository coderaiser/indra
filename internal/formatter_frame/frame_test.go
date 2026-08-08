package formatter_frame_test

import (
	"regexp"
	"strings"
	"testing"

	formframe "coderaiser/indra/internal/formatter_frame"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var ansiRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

func strip(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

var src = []byte("package p\n\nfunc f() {\n\tx := 1\n\t_ = 0\n}\n")
var place1 = types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 4, Column: 2}}

func TestFrame(t *testing.T) {
	Test(t, "frame: mid-run returns empty (percent on stderr)", func(t *T) {
		out := formframe.Format("a.go", src, nil, 0, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})

	Test(t, "frame: last file returns codeframe output", func(t *T) {
		out := formframe.Format("a.go", src, []types.Place{place1}, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "frame: last file clean returns empty", func(t *T) {
		out := formframe.Format("a.go", src, nil, 4, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})

	Test(t, "frame: single file returns codeframe immediately", func(t *T) {
		out := formframe.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "> 4"))
		t.End()
	})
}
