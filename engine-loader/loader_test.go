package engine_loader

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
)

func replacerFuncs(name, path string) PluginFuncs {
	return PluginFuncs{
		Name:    name,
		Path:    path,
		Report:  func() string { return "r:" + name },
		Match:   func() types.Matcher { return types.Matcher{"a": nil} },
		Replace: func() types.Replacer { return types.Replacer{"a": "b"} },
	}
}

func traverserFuncs(name, path string) PluginFuncs {
	return PluginFuncs{
		Name:     name,
		Path:     path,
		Report:   func(node ast.Node) string { return "t:" + name },
		Traverse: func() types.Traverser { return types.Traverser{"*ast.File": fileVisitor} },
		Fix:      func(node ast.Node, opts map[string]any) {},
	}
}

func fileVisitor(node ast.Node, push func(ast.Node)) {}

func names(kinds []PluginKind) map[string]bool {
	m := make(map[string]bool, len(kinds))
	for _, k := range kinds {
		m[k.Name()] = true
	}
	return m
}

func TestLoadAllEnabled(t *testing.T) {
	plugins := []PluginFuncs{
		replacerFuncs("remove-skip", "p/remove-skip"),
		traverserFuncs("remove-unused-import", "p/remove-unused-import"),
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
	plugins := []PluginFuncs{replacerFuncs("skip", "p/skip")}
	cfg := Config{"skip": {Enabled: false}}
	got := Load(plugins, cfg)
	if len(got) != 0 {
		t.Fatalf("expected 0 plugins, got %d", len(got))
	}
}

// nestedFuncs builds a nested (grouping) plugin whose Rules reference path.
func nestedFuncs(name, path string, rules types.Nested) PluginFuncs {
	return PluginFuncs{Name: name, Path: path, Rules: rules}
}

func TestLoadPrefixDisabled(t *testing.T) {
	// Both remove-skip and add-t-end exist as top-level rules AND as tape/*
	// sub-rules. Disabling the "tape" prefix keeps only the top-level rules.
	plugins := []PluginFuncs{
		replacerFuncs("remove-skip", "p/remove-skip"),
		replacerFuncs("add-t-end", "p/add-t-end"),
		nestedFuncs("tape", "p/tape", types.Nested{
			"remove-skip": "p/remove-skip",
			"add-t-end":   "p/add-t-end",
		}),
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

func TestLoadOffRespected(t *testing.T) {
	// The tape/remove-skip sub-rule is Off() by default; only the top-level
	// remove-skip rule survives.
	plugins := []PluginFuncs{
		replacerFuncs("remove-skip", "p/remove-skip"),
		nestedFuncs("tape", "p/tape", types.Nested{"remove-skip": types.Off("p/remove-skip")}),
	}
	got := Load(plugins, Config{})
	nms := names(got)
	if nms["tape/remove-skip"] {
		t.Fatalf("expected Off() tape/remove-skip disabled, got %v", nms)
	}
	if !nms["remove-skip"] {
		t.Fatalf("expected top-level remove-skip enabled, got %v", nms)
	}
}

func TestLoadOffOverriddenOn(t *testing.T) {
	// Config turning tape/remove-skip on overtakes the default Off().
	plugins := []PluginFuncs{
		replacerFuncs("remove-skip", "p/remove-skip"),
		nestedFuncs("tape", "p/tape", types.Nested{"remove-skip": types.Off("p/remove-skip")}),
	}
	cfg := Config{"tape/remove-skip": {Enabled: true}}
	got := Load(plugins, cfg)
	nms := names(got)
	if !nms["tape/remove-skip"] {
		t.Fatalf("expected tape/remove-skip enabled by config, got %v", nms)
	}
}

func TestLoadMissingNestedPath(t *testing.T) {
	// nested path with no matching PluginFuncs is skipped
	plugins := []PluginFuncs{
		replacerFuncs("x", "p/x"),
		nestedFuncs("tape", "p/tape", types.Nested{"ghost": "p/ghost"}),
	}
	got := Load(plugins, Config{})
	if len(got) != 1 {
		t.Fatalf("expected only top-level x, got %d", len(got))
	}
	if got[0].Name() != "x" {
		t.Fatalf("expected x, got %q", got[0].Name())
	}
}

func TestLoadWrongSignaturePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for wrong signature")
		}
	}()
	bad := PluginFuncs{Name: "bad", Report: "not-a-func"}
	Load([]PluginFuncs{bad}, Config{})
}

func TestLoadUnknownKindPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown kind")
		}
	}()
	bad := PluginFuncs{
		Name:   "bad",
		Report: func() string { return "x" },
		Match:  func() types.Matcher { return nil },
		// no Replace, no Traverse/Fix
	}
	Load([]PluginFuncs{bad}, Config{})
}

func TestDefaultConfigEmpty(t *testing.T) {
	if len(DefaultConfig()) != 0 {
		t.Fatal("expected empty default config")
	}
}

// TestReplacerPluginAccessors exercises the resolved replacer's methods.
func TestReplacerPluginAccessors(t *testing.T) {
	p := replacerFuncs("rp", "p/rp")
	k := resolveKind(p, "rp")
	rp, ok := k.(ReplacerPlugin)
	if !ok {
		t.Fatalf("expected ReplacerPlugin, got %T", k)
	}
	rp.pluginKind()
	if rp.Name() != "rp" || rp.Report() != "r:rp" {
		t.Fatalf("unexpected name/report: %q %q", rp.Name(), rp.Report())
	}
	if rp.Match()["a"] != nil || rp.Replace()["a"] != "b" {
		t.Fatal("unexpected Match/Replace accessors")
	}
}

// TestTraverserPluginAccessors exercises the resolved traverser's methods.
func TestTraverserPluginAccessors(t *testing.T) {
	p := traverserFuncs("tp", "p/tp")
	k := resolveKind(p, "tp")
	tp, ok := k.(TraverserPlugin)
	if !ok {
		t.Fatalf("expected TraverserPlugin, got %T", k)
	}
	tp.pluginKind()
	if tp.Name() != "tp" || tp.Report(nil) != "t:tp" {
		t.Fatalf("unexpected name/report: %q %q", tp.Name(), tp.Report(nil))
	}
	if tp.Traverse()["*ast.File"] == nil {
		t.Fatal("expected Traverse accessor to return visitor")
	}
	tp.Fix(nil, nil) // no-op, just exercises the wrapper
}

// TestMissingReportPanics covers invokeReport's nil-Report guard.
func TestMissingReportPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing Report")
		}
	}()
	Load([]PluginFuncs{{Name: "bad", Match: func() types.Matcher { return nil }, Replace: func() types.Replacer { return nil }}}, Config{})
}

// TestNothingReturnPanics covers funcValue's "not a func" guard.
func TestNothingReturnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for non-func field")
		}
	}()
	Load([]PluginFuncs{{Name: "bad", Report: func() string { return "x" }, Match: "no", Replace: "no"}}, Config{})
}

// TestFixWrongShapePanics covers funcValue's Fix two-arg guard with a bad shape.
func TestFixWrongShapePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for wrong Fix shape")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Report:   func() string { return "x" },
		Traverse: func() types.Traverser { return nil },
		Fix:      func() {},
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestUnknownFieldPanics covers fieldOf's default branch.
func TestUnknownFieldPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown field")
		}
	}()
	fieldOf(PluginFuncs{Name: "p"}, "Yolo")
}

// TestEntryPathUnknown covers entryPath's default branch.
func TestEntryPathUnknown(t *testing.T) {
	if entryPath(42) != "" {
		t.Fatal("expected empty path for unknown value type")
	}
}

// TestEntryPathPluginEntry covers entryPath's PluginEntry branch.
func TestEntryPathPluginEntry(t *testing.T) {
	if entryPath(types.PluginEntry{Path: "x", Enabled: true}) != "x" {
		t.Fatal("expected PluginEntry path")
	}
}

// TestEntryPathString covers entryPath's string branch.
func TestEntryPathString(t *testing.T) {
	if entryPath("y") != "y" {
		t.Fatal("expected string path")
	}
}

// TestMultiReturnPanics covers funcValue's zero-arg multi-return guard.
func TestMultiReturnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for multi-return func")
		}
	}()
	bad := PluginFuncs{
		Name:    "bad",
		Report:  func() (string, error) { return "", nil },
		Match:   func() types.Matcher { return nil },
		Replace: func() types.Replacer { return nil },
	}
	Load([]PluginFuncs{bad}, Config{})
}


// TestResolveKindReplacerWithoutMatch covers a replacer whose Match is nil.
func TestResolveKindReplacerWithoutMatch(t *testing.T) {
	pf := PluginFuncs{
		Name:    "no-match",
		Path:    "p/no-match",
		Report:  func() string { return "r" },
		Replace: func() types.Replacer { return types.Replacer{"a": "b"} },
		// Match intentionally nil
	}
	kinds := Load([]PluginFuncs{pf}, Config{})
	if len(kinds) != 1 {
		t.Fatalf("expected 1 kind, got %d", len(kinds))
	}
	if _, ok := kinds[0].(ReplacerPlugin); !ok {
		t.Fatalf("expected ReplacerPlugin, got %T", kinds[0])
	}
}

// TestReplacerPluginMatchNilReturnsEmpty covers a replacer without Match
// returning an empty Matcher.
func TestReplacerPluginMatchNilReturnsEmpty(t *testing.T) {
	pf := PluginFuncs{
		Name:    "no-match",
		Path:    "p/no-match",
		Report:  func() string { return "r" },
		Replace: func() types.Replacer { return types.Replacer{"a": "b"} },
	}
	kinds := Load([]PluginFuncs{pf}, Config{})
	rp := kinds[0].(ReplacerPlugin)
	result := rp.Match()
	if len(result) != 0 {
		t.Fatalf("expected empty Matcher, got %d entries", len(result))
	}
}

// TestMissingTraverserReportPanics covers invokeTraverserReport's nil guard.
func TestMissingTraverserReportPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing Report")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Traverse: func() types.Traverser { return nil },
		Fix:      func(node ast.Node, opts map[string]any) {},
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestTraverserReportWrongReturnPanics covers funcValue's Report single-return
// guard for a traverser Report that returns two values.
func TestTraverserReportWrongReturnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for wrong Report return count")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Report:   func(node ast.Node) (string, error) { return "", nil },
		Traverse: func() types.Traverser { return nil },
		Fix:      func(node ast.Node, opts map[string]any) {},
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestTraverserReportTooManyArgsPanics covers funcValue's Report arity guard
// for a traverser Report that takes two arguments.
func TestTraverserReportTooManyArgsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for Report with two args")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Report:   func(a, b ast.Node) string { return "" },
		Traverse: func() types.Traverser { return nil },
		Fix:      func(node ast.Node, opts map[string]any) {},
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestFixWrongOutputPanics covers funcValue's Fix output-count guard when Fix
// has the right two inputs but returns a value.
func TestFixWrongOutputPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for Fix with a return value")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Report:   func(node ast.Node) string { return "x" },
		Traverse: func() types.Traverser { return nil },
		Fix:      func(node ast.Node, opts map[string]any) string { return "" },
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestTraverseMultiReturnPanics covers funcValue's default guard for a
// Traverse field that returns multiple values.
func TestTraverseMultiReturnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for multi-return Traverse")
		}
	}()
	bad := PluginFuncs{
		Name:     "bad",
		Report:   func(node ast.Node) string { return "x" },
		Traverse: func() (types.Traverser, error) { return nil, nil },
		Fix:      func(node ast.Node, opts map[string]any) {},
	}
	Load([]PluginFuncs{bad}, Config{})
}

// TestMatchWrongReturnPanics covers mustFunc's type-assertion guard when Match
// is a func with the right arity but the wrong return type.
func TestMatchWrongReturnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for wrong Match return type")
		}
	}()
	bad := PluginFuncs{
		Name:    "bad",
		Report:  func() string { return "x" },
		Match:   func() int { return 1 },
		Replace: func() types.Replacer { return nil },
	}
	Load([]PluginFuncs{bad}, Config{})
}

