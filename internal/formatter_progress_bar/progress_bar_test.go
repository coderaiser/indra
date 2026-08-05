package formatter_progress_bar_test

import (
	"strings"
	"testing"

	pb "coderaiser/indra/internal/formatter_progress_bar"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var places1 = []types.Place{
	{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 5, Column: 2}},
}

func TestProgressBarMidRunReturnsEmpty(t *testing.T) {
	Test(t, "progress-bar: mid-run returns empty string", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", nil, 0, 10, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarLastFileNoIssues(t *testing.T) {
	Test(t, "progress-bar: last file no issues returns empty", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", nil, 4, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarLastFileWithIssues(t *testing.T) {
	Test(t, "progress-bar: last file with issues returns dump output", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", places1, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})
}

func TestProgressBarForceShow(t *testing.T) {
	Test(t, "progress-bar: INDRA_PROGRESS_BAR=1 forces show", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")
		out := pb.Format("foo.go", nil, 0, 10, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarShouldShowForcedOff(t *testing.T) {
	Test(t, "progress-bar: ShouldShow forced off by env", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		t.Ok(!pb.ShouldShow(1000))
		t.End()
	})
}

func TestRenderBarContainsBlocks(t *testing.T) {
	Test(t, "progress-bar: RenderBar contains block chars", func(t *T) {
		result := pb.RenderBar(5, 10, "#6fbdf1")
		t.Ok(strings.Contains(result, "█") || strings.Contains(result, "░"))
		t.End()
	})
}

func TestRenderBarZeroTotal(t *testing.T) {
	Test(t, "progress-bar: RenderBar zero total returns empty bar", func(t *T) {
		result := pb.RenderBar(0, 0, "#6fbdf1")
		t.Ok(strings.Contains(result, "░"))
		t.End()
	})
}

func TestTermWidthEnv(t *testing.T) {
	Test(t, "progress-bar: TermWidth respects env var", func(t *T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "120")
		result := pb.TermWidth()
		t.Equal(result, 120)

		t.End()
	})
}

func TestTermWidthDefault(t *testing.T) {
	Test(t, "progress-bar: TermWidth returns positive default", func(t *T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "")
		t.Ok(pb.TermWidth() > 0)
		t.End()
	})
}

func TestProgressBarShowTruncatedLine(t *testing.T) {
	Test(t, "progress-bar: forced show truncates overlong line", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "5")
		out := pb.Format("a-very-long-filename-indeed.go", nil, 2, 3, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarShowLastFileReturnsDump(t *testing.T) {
	Test(t, "progress-bar: forced show on last file returns dump output", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")
		out := pb.Format("foo.go", places1, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})
}

func TestProgressBarShowMidRunReturnsEmpty(t *testing.T) {
	Test(t, "progress-bar: forced show mid-run returns empty", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")
		out := pb.Format("foo.go", nil, 0, 10, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarShowWithErrors(t *testing.T) {
	Test(t, "progress-bar: forced show with errors renders count", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")
		out := pb.Format("foo.go", nil, 3, 10, 1, 2)
		t.Equal(out, "")
		t.End()
	})
}

func TestRenderBarClampsFilled(t *testing.T) {
	Test(t, "progress-bar: RenderBar clamps filled to bar width", func(t *T) {
		result := pb.RenderBar(1000, 10, "#6fbdf1")
		t.Ok(strings.Contains(result, "█"))
		t.End()
	})
}

func TestRenderBarNonHexColor(t *testing.T) {
	Test(t, "progress-bar: RenderBar passes through non hex color", func(t *T) {
		result := pb.RenderBar(5, 10, "red")
		t.Ok(strings.Contains(result, "red"))
		t.End()
	})
}

func TestTruncateLongString(t *testing.T) {
	Test(t, "progress-bar: Truncate shortens long string", func(t *T) {
		result := pb.Truncate("abcdefghij", 5)
		t.Ok(strings.Contains(result, "..."))
		t.End()
	})
}

func TestTruncateANSIWithinLen(t *testing.T) {
	Test(t, "progress-bar: TruncateANSI returns unchanged when within length", func(t *T) {
		result := pb.TruncateANSI("\x1b[31mfoo\x1b[0m", 10)
		t.Equal(result, "\x1b[31mfoo\x1b[0m")
		t.End()
	})
}

func TestTruncateANSIOverLen(t *testing.T) {
	Test(t, "progress-bar: TruncateANSI shortens visible chars", func(t *T) {
		result := pb.TruncateANSI("\x1b[31mabcdefghij\x1b[0m", 3)
		t.Ok(strings.Contains(result, "abc"))
		t.End()
	})
}

func TestHexToANSIUpper(t *testing.T) {
	Test(t, "progress-bar: hexToANSI handles uppercase hex", func(t *T) {
		result := pb.RenderBar(1, 2, "#6FBDF1")
		t.Ok(strings.Contains(result, "\x1b[38;2;"))
		t.End()
	})
}

func TestHexToANSIInvalid(t *testing.T) {
	Test(t, "progress-bar: hexToANSI returns invalid color unchanged", func(t *T) {
		result := pb.RenderBar(1, 2, "nocolor")
		t.Ok(strings.HasPrefix(result, "nocolor"))
		t.End()
	})
}

func TestShouldShowInvalidMinFallback(t *testing.T) {
	Test(t, "progress-bar: ShouldShow ignores invalid min env", func(t *T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		t.Ok(!pb.ShouldShow(1))
		t.End()
	})
}
