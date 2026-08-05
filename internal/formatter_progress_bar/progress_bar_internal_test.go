package formatter_progress_bar

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestConfigureColor(t *testing.T) {
	tape.Test(t, "Configure: sets color used by RenderBar", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: "#ff0000", MinCount: 0})
		t.Ok(strings.Contains(RenderBar(1, 1, cfg.Color), "255;0;0"))
		t.End()
	})
}

func TestConfigureEmptyColorKeepsDefault(t *testing.T) {
	tape.Test(t, "Configure: empty color keeps defaultColor", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: "", MinCount: 0})
		t.Equal(cfg.Color, defaultColor)
		t.End()
	})
}

func TestConfigureMinBelowThreshold(t *testing.T) {
	tape.Test(t, "Configure: MinCount below threshold hides bar", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: defaultColor, MinCount: 5})
		t.Equal(ShouldShow(4), false)
		t.End()
	})
}

func TestConfigureMinAtThreshold(t *testing.T) {
	tape.Test(t, "Configure: MinCount at threshold shows bar", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		Configure(Config{Color: defaultColor, MinCount: 5})
		t.Equal(ShouldShow(5), true)
		t.End()
	})
}

func TestShouldShowEnv1Overrides(t *testing.T) {
	tape.Test(t, "ShouldShow: INDRA_PROGRESS_BAR=1 overrides MinCount", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		t.TB().Setenv("INDRA_PROGRESS_BAR", "1")
		Configure(Config{Color: defaultColor, MinCount: 999})
		t.Equal(ShouldShow(1), true)
		t.End()
	})
}

func TestShouldShowEnv0Overrides(t *testing.T) {
	tape.Test(t, "ShouldShow: INDRA_PROGRESS_BAR=0 overrides MinCount", func(t *tape.T) {
		cfg = Config{Color: defaultColor, MinCount: 0}
		t.TB().Setenv("INDRA_PROGRESS_BAR", "0")
		Configure(Config{Color: defaultColor, MinCount: 0})
		t.Equal(ShouldShow(999), false)
		t.End()
	})
}
