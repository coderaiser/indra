package formatter_progress_bar

import (
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestConfigureColor(t *testing.T) {
	Test(t, "Configure: sets color used by RenderBar", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: "#ff0000", MinCount: 0})
		t.Ok(strings.Contains(RenderBar(1, 1, cfg.Color), "255;0;0"))
		t.End()
	})
}

func TestConfigureEmptyColorKeepsDefault(t *testing.T) {
	Test(t, "Configure: empty color keeps defaultColor", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: "", MinCount: 0})
		t.Equal(cfg.Color, defaultColor)
		t.End()
	})
}

func TestConfigureMinBelowThreshold(t *testing.T) {
	Test(t, "Configure: MinCount below threshold hides bar", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: defaultColor, MinCount: 5})
		result := ShouldShow(4)
		t.Equal(result, false)

		t.End()
	})
}

func TestConfigureMinAtThreshold(t *testing.T) {
	Test(t, "Configure: MinCount at threshold shows bar", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: defaultColor, MinCount: 5})
		result := ShouldShow(5)
		t.Equal(result, true)

		t.End()
	})
}

func TestShouldShowEnv1Overrides(t *testing.T) {
	Test(t, "ShouldShow: INDRA_PROGRESS_BAR=1 overrides MinCount", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		Configure(Config{Color: defaultColor, MinCount: 999})
		result := ShouldShow(1)
		t.Equal(result, true)

		t.End()
	})
}

func TestShouldShowEnv0Overrides(t *testing.T) {
	Test(t, "ShouldShow: INDRA_PROGRESS_BAR=0 overrides MinCount", func(t *T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		Configure(Config{Color: defaultColor, MinCount: 0})
		result := ShouldShow(999)
		t.Equal(result, false)

		t.End()
	})
}
