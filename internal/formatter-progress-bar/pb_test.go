package formatter_progress_bar_test

import (
	"strings"
	"testing"

	pb "coderaiser/indra/internal/formatter-progress-bar"
	"coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

var places1 = []types.Place{
	{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 5, Column: 2}},
}

func TestProgressBarMidRunReturnsEmpty(t *testing.T) {
	tape.Test(t, "progress-bar: mid-run returns empty string", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", nil, 0, 10, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarLastFileNoIssues(t *testing.T) {
	tape.Test(t, "progress-bar: last file no issues returns empty", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", nil, 4, 5, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarLastFileWithIssues(t *testing.T) {
	tape.Test(t, "progress-bar: last file with issues returns dump output", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		out := pb.Format("foo.go", places1, 4, 5, 1, 1)
		t.Ok(strings.Contains(out, "1 error in 1 file"))
		t.End()
	})
}

func TestProgressBarForceShow(t *testing.T) {
	tape.Test(t, "progress-bar: INDRA_PROGRESS_BAR=1 forces show", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		t.TB().Setenv("INDRA_TERM_WIDTH", "80")
		out := pb.Format("foo.go", nil, 0, 10, 0, 0)
		t.Equal(out, "")
		t.End()
	})
}

func TestProgressBarShouldShowTrue(t *testing.T) {
	tape.Test(t, "progress-bar: ShouldShow true when count >= min", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR_MIN", "3")
		t.Ok(pb.ShouldShow(3))
		t.End()
	})
}

func TestProgressBarShouldShowFalse(t *testing.T) {
	tape.Test(t, "progress-bar: ShouldShow false when count < min", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "")
		t.TB().Setenv("INDRA_PROGRESS_BAR_MIN", "10")
		t.Ok(!pb.ShouldShow(3))
		t.End()
	})
}

func TestProgressBarShouldShowForcedOff(t *testing.T) {
	tape.Test(t, "progress-bar: ShouldShow forced off by env", func(t *tape.T) {
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		t.Ok(!pb.ShouldShow(1000))
		t.End()
	})
}

func TestRenderBarContainsBlocks(t *testing.T) {
	tape.Test(t, "progress-bar: RenderBar contains block chars", func(t *tape.T) {
		result := pb.RenderBar(5, 10, "#6fbdf1")
		t.Ok(strings.Contains(result, "█") || strings.Contains(result, "░"))
		t.End()
	})
}

func TestRenderBarZeroTotal(t *testing.T) {
	tape.Test(t, "progress-bar: RenderBar zero total returns empty bar", func(t *tape.T) {
		result := pb.RenderBar(0, 0, "#6fbdf1")
		t.Ok(strings.Contains(result, "░"))
		t.End()
	})
}

func TestTermWidthEnv(t *testing.T) {
	tape.Test(t, "progress-bar: TermWidth respects env var", func(t *tape.T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "120")
		t.Equal(pb.TermWidth(), 120)
		t.End()
	})
}

func TestTermWidthDefault(t *testing.T) {
	tape.Test(t, "progress-bar: TermWidth returns positive default", func(t *tape.T) {
		t.TB().Setenv("INDRA_TERM_WIDTH", "")
		t.Ok(pb.TermWidth() > 0)
		t.End()
	})
}
