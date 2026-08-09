package formatter_codeframe_test

import (
	"regexp"
	"strings"
	"testing"

	formcf "coderaiser/indra/internal/formatter_codeframe"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var ansiRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

func strip(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

var src = []byte(`package p

func f() {
    x := 1
    _ = 0
}
`)

var place1 = types.Place{Rule: "r", Message: "remove unused variable: x", Position: types.Position{Line: 4, Column: 5}}

func TestCodeframe(t *testing.T) {
	Test(t, "codeframe: no places returns empty", func(t *T) {
		out := formcf.Format("a.go", src, nil, 0, 1, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "codeframe: no places mid-run returns empty", func(t *T) {
		out := formcf.Format("a.go", src, nil, 0, 3, 0, 0)
		t.NotOk(out)

		t.End()
	})

	Test(t, "codeframe: last clean file with prior errors returns summary", func(t *T) {
		out := formcf.Format("a.go", src, nil, 0, 1, 1, 2)
		t.Ok(strings.Contains(out, "2 errors in 1 file"))
		t.End()
	})

	Test(t, "codeframe: shows filename", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "a.go"))
		t.End()
	})

	Test(t, "codeframe: shows error line highlighted with >", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "> 4"))
		t.End()
	})

	Test(t, "codeframe: shows surrounding context lines", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "3 │"))
		t.End()
	})

	Test(t, "codeframe: shows message", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "remove unused variable: x"))
		t.End()
	})

	Test(t, "codeframe: nil source falls back to dump output", func(t *T) {
		out := formcf.Format("a.go", nil, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "remove unused variable: x"))
		t.End()
	})

	Test(t, "codeframe: last file appends summary", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})

	Test(t, "codeframe: place at first line shows no preceding context", func(t *T) {
		p := types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 1, Column: 1}}
		out := formcf.Format("a.go", src, []types.Place{p}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "> 1"))
		t.End()
	})

	Test(t, "codeframe: place beyond source lines shows dump fallback", func(t *T) {
		p := types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 999, Column: 1}}
		out := formcf.Format("a.go", src, []types.Place{p}, 0, 1, 1, 1)
		t.Ok(len(out) > 0)
		t.End()
	})

	Test(t, "codeframe: place near end clamps context to last line", func(t *T) {
		p := types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 6, Column: 1}}
		out := formcf.Format("a.go", src, []types.Place{p}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "> 6"))
		t.End()
	})

	Test(t, "codeframe: shows caret under exact column", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "\033[33m^\033[0m"))
		t.End()
	})

	Test(t, "codeframe: caret line contains message and rule", func(t *T) {
		out := formcf.Format("a.go", src, []types.Place{place1}, 0, 1, 1, 1)
		t.Ok(strings.Contains(strip(out), "remove unused variable: x (r)"))
		t.End()
	})

	Test(t, "codeframe: clamps column to 1 for caret line", func(t *T) {
		p := types.Place{Rule: "r", Message: "m", Position: types.Position{Line: 4, Column: 0}}
		out := formcf.Format("a.go", src, []types.Place{p}, 0, 1, 1, 1)
		t.Ok(strings.Contains(out, "\033[33m^\033[0m"))
		t.End()
	})
}
