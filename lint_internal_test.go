package indra

import (
	"testing"

	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/internal/config"

	. "github.com/coderaiser/go-tape"
)

// pluginItems builds runner.PluginItem with the given rule names.
func pluginItems(names ...string) []runner.PluginItem {
	out := make([]runner.PluginItem, 0, len(names))
	for _, n := range names {
		out = append(out, runner.PluginItem{Rule: n})
	}
	return out
}

func TestFilterPluginsExact(t *testing.T) {
	Test(t, "filterPlugins: keeps exact rule match", func(t *T) {
		out := filterPlugins(pluginItems("a", "b"), []string{"b"})
		result := len(out)
		t.Equal(result, 1)

		t.End()
	})
}

func TestFilterPluginsGroup(t *testing.T) {
	Test(t, "filterPlugins: keeps group prefix matches", func(t *T) {
		out := filterPlugins(pluginItems("tape/a", "tape/b", "x"), []string{"tape"})
		result := len(out)
		t.Equal(result, 2)

		t.End()
	})
}

func TestFilterPluginsDropsOthers(t *testing.T) {
	Test(t, "filterPlugins: drops non matching rules", func(t *T) {
		out := filterPlugins(pluginItems("a", "b"), []string{"grp"})
		result := len(out)
		t.Equal(result, 0)

		t.End()
	})
}

func TestConfigForFileNoOverride(t *testing.T) {
	Test(t, "configForFile: keeps base rules with no match", func(t *T) {
		cfg := config.Config{Rules: map[string]string{"tape": "on"}}
		lc := configForFile(cfg, "plain.go")
		t.Ok(lc["tape"].Enabled)

		t.End()
	})
}

func TestConfigForFileOverrideUnmatched(t *testing.T) {
	Test(t, "configForFile: ignored when pattern does not match file", func(t *T) {
		cfg := config.Config{
			Rules: map[string]string{"tape": "on"},
			Match: config.MatchConfig{"skip_*.go": {"tape": "off"}},
		}
		lc := configForFile(cfg, "plain.go")
		t.Ok(lc["tape"].Enabled)

		t.End()
	})
}

func TestConfigForFileOverrideOff(t *testing.T) {
	Test(t, "configForFile: match override switches a rule off", func(t *T) {
		cfg := config.Config{
			Rules: map[string]string{"tape": "on"},
			Match: config.MatchConfig{"skip_*.go": {"tape": "off"}},
		}
		lc := configForFile(cfg, "skip_a.go")
		t.NotOk(lc["tape"].Enabled)

		t.End()
	})
}

func TestConfigForFileOverrideOn(t *testing.T) {
	Test(t, "configForFile: match override switches a rule on", func(t *T) {
		cfg := config.Config{
			Rules: map[string]string{"tape": "off"},
			Match: config.MatchConfig{"skip_*.go": {"tape": "on"}},
		}
		lc := configForFile(cfg, "skip_a.go")
		t.Ok(lc["tape"].Enabled)

		t.End()
	})
}
