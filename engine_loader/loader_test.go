package engine_loader

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
)

// testReplacer is a minimal replacer plugin: Replace but no Match.
type testReplacer struct{}

func (testReplacer) Report() string          { return "r" }
func (testReplacer) Replace() types.Replacer { return types.Replacer{"a": "b"} }

// testReplacerMatch is a replacer that also exposes a Match guard.
type testReplacerMatch struct{}

func (testReplacerMatch) Report() string          { return "rm" }
func (testReplacerMatch) Match() types.Matcher    { return types.Matcher{"p": nil} }
func (testReplacerMatch) Replace() types.Replacer { return types.Replacer{"p": "q"} }

// testTraverser is a minimal traverser plugin.
type testTraverser struct{}

func (testTraverser) Report(_ types.Path) string         { return "t" }
func (testTraverser) Traverse() types.Traverser          { return types.Traverser{"*ast.File": fileVisitor} }
func (testTraverser) Fix(_ types.Path, _ map[string]any) {}

func fileVisitor(p types.Path, push func(types.Path)) {}

// ── malformed shapes ─────────────────────────────────────────────────────────

// badReplaceSignature has a Replace method of the wrong type.
type badReplaceSignature struct{}

func (badReplaceSignature) Report() string { return "b" }
func (badReplaceSignature) Replace() int   { return 1 }

// badTraverseSignature has a Traverse method of the wrong type.
type badTraverseSignature struct{}

func (badTraverseSignature) Report(_ ast.Node) string { return "b" }
func (badTraverseSignature) Traverse() (types.Traverser, error) {
	return nil, nil
}
func (badTraverseSignature) Fix(_ ast.Node, _ map[string]any) {}

// badMatchSignature has a valid Replace but a Match of the wrong type.
type badMatchSignature struct{}

func (badMatchSignature) Report() string          { return "b" }
func (badMatchSignature) Match() int              { return 1 }
func (badMatchSignature) Replace() types.Replacer { return types.Replacer{} }

// noReport is a replacer without a Report method.
type noReport struct{}

func (noReport) Replace() types.Replacer { return types.Replacer{} }

// noReportTraverser is a traverser without a Report method.
type noReportTraverser struct{}

func (noReportTraverser) Traverse() types.Traverser        { return types.Traverser{} }
func (noReportTraverser) Fix(_ ast.Node, _ map[string]any) {}

// noFixTraverser is a traverser without a Fix method.
type noFixTraverser struct{}

func (noFixTraverser) Report(_ ast.Node) string  { return "t" }
func (noFixTraverser) Traverse() types.Traverser { return types.Traverser{} }

// emptyPlugin carries no shape methods at all.
type emptyPlugin struct{}

func names(kinds []PluginKind) map[string]bool {
	m := make(map[string]bool, len(kinds))
	for _, k := range kinds {
		m[k.Name()] = true
	}
	return m
}

func TestLoadAllEnabled(t *testing.T) {
	plugins := []PluginFuncs{
		{Name: "remove-skip", Plugin: testReplacer{}},
		{Name: "remove-unused-import", Plugin: testTraverser{}},
	}
	got := Load(plugins, Config{})
	if len(got) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(got))
	}
	if _, ok := got[0].(ReplacerPlugin); !ok {
		t.Fatalf("expected ReplacerPlugin for first, got %T", got[0])
	}
	if _, ok := got[1].(TraverserPlugin); !ok {
		t.Fatalf("expected TraverserPlugin for second, got %T", got[1])
	}
}

func TestLoadExactDisabled(t *testing.T) {
	plugins := []PluginFuncs{{Name: "skip", Plugin: testReplacer{}}}
	cfg := Config{"skip": {Enabled: false}}
	got := Load(plugins, cfg)
	if len(got) != 0 {
		t.Fatalf("expected 0 plugins, got %d", len(got))
	}
}

func TestLoadGroupExpands(t *testing.T) {
	plugins := []PluginFuncs{{
		Name: "tape",
		Rules: []types.Rule{
			{Name: "remove-skip", Plugin: testReplacer{}},
			{Name: "add-t-end", Plugin: testReplacerMatch{}},
		},
	}}
	got := Load(plugins, Config{})
	nms := names(got)
	if !nms["tape/remove-skip"] || !nms["tape/add-t-end"] {
		t.Fatalf("expected tape/* rules, got %v", nms)
	}
}

func TestLoadPrefixDisabled(t *testing.T) {
	// A standalone leaf and a group sharing the same rule names. Disabling the
	// "tape" prefix keeps only the standalone rules.
	plugins := []PluginFuncs{
		{Name: "remove-skip", Plugin: testReplacer{}},
		{Name: "add-t-end", Plugin: testReplacer{}},
		{Name: "tape", Rules: []types.Rule{
			{Name: "remove-skip", Plugin: testReplacer{}},
			{Name: "add-t-end", Plugin: testReplacer{}},
		}},
	}
	cfg := Config{"tape": {Enabled: false}}
	got := Load(plugins, cfg)
	nms := names(got)
	if !nms["remove-skip"] || !nms["add-t-end"] {
		t.Fatalf("expected top-level rules kept, got %v", nms)
	}
	if nms["tape/remove-skip"] || nms["tape/add-t-end"] {
		t.Fatalf("expected tape/* rules disabled, got %v", nms)
	}
}

func TestLoadGroupRuleOffByConfig(t *testing.T) {
	// A group rule is on by default; config "off" disables it.
	plugins := []PluginFuncs{{
		Name: "indra",
		Rules: []types.Rule{
			{Name: "remove-useless-match", Plugin: testTraverser{}},
			{Name: "other", Plugin: testReplacer{}},
		},
	}}
	cfg := Config{"indra/remove-useless-match": {Enabled: false}}
	got := Load(plugins, cfg)
	nms := names(got)
	if nms["indra/remove-useless-match"] {
		t.Fatalf("expected rule off by config, got %v", nms)
	}
	if !nms["indra/other"] {
		t.Fatalf("expected other rule present, got %v", nms)
	}
}

func TestResolveReplacerAccessors(t *testing.T) {
	k := resolve(testReplacer{}, "rp", "rp")
	rp, ok := k.(ReplacerPlugin)
	if !ok {
		t.Fatalf("expected ReplacerPlugin, got %T", k)
	}
	rp.pluginKind()
	if rp.Name() != "rp" || rp.Report() != "r" {
		t.Fatalf("unexpected name/report: %q %q", rp.Name(), rp.Report())
	}
	if len(rp.Match()) != 0 || rp.Replace()["a"] != "b" {
		t.Fatal("unexpected Match/Replace accessors")
	}
}

func TestResolveReplacerWithMatch(t *testing.T) {
	k := resolve(testReplacerMatch{}, "rm", "rm")
	rp := k.(ReplacerPlugin)
	if rp.Match()["p"] != nil || rp.Replace()["p"] != "q" {
		t.Fatal("unexpected Match/Replace values")
	}
}

func TestResolveTraverserAccessors(t *testing.T) {
	k := resolve(testTraverser{}, "tp", "tp")
	tp, ok := k.(TraverserPlugin)
	if !ok {
		t.Fatalf("expected TraverserPlugin, got %T", k)
	}
	tp.pluginKind()
	if tp.Name() != "tp" || tp.Report(types.Path{}) != "t" {
		t.Fatalf("unexpected name/report: %q %q", tp.Name(), tp.Report(types.Path{}))
	}
	if tp.Traverse()["*ast.File"] == nil {
		t.Fatal("expected Traverse accessor to return visitor")
	}
	tp.Fix(types.Path{}, nil) // no-op, exercises the wrapper
}

func TestDefaultConfigEmpty(t *testing.T) {
	if len(DefaultConfig()) != 0 {
		t.Fatal("expected empty default config")
	}
}

// ── panic / malformed-shape paths ────────────────────────────────────────────

func catchPanic(fn func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			s, _ := r.(string)
			msg = s
		}
	}()
	fn()
	return msg
}

func TestUnknownKindPanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: emptyPlugin{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for unknown kind")
	}
}

func TestNilPluginPanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x"}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for nil plugin")
	}
}

func TestMissingReplacerReportPanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: noReport{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for missing replacer Report")
	}
}

func TestMissingTraverserReportPanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: noReportTraverser{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for missing traverser Report")
	}
}

func TestMissingFixPanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: noFixTraverser{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for missing Fix")
	}
}

func TestReplaceWrongSignaturePanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: badReplaceSignature{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for wrong Replace signature")
	}
}

func TestTraverseWrongSignaturePanics(t *testing.T) {
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: badTraverseSignature{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for wrong Traverse signature")
	}
}

func TestMatchWrongSignaturePanics(t *testing.T) {
	// A struct with an invalid Match alongside a valid Replace must panic.
	if msg := catchPanic(func() { Load([]PluginFuncs{{Name: "x", Plugin: badMatchSignature{}}}, Config{}) }); msg == "" {
		t.Fatal("expected panic for wrong Match signature")
	}
}
