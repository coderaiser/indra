// Package engine_loader resolves plugin packages into runnable plugin kinds
// and filters them by rule configuration.
package engine_loader

import (
	"go/ast"
	"reflect"

	"coderaiser/indra/types"
)

// PluginFuncs is a set of exported funcs from one plugin package.
// It is passed by plugins.go which imports the packages statically.
// A nested plugin carries Rules instead of Report/Match/Replace/Traverse/Fix.
type PluginFuncs struct {
	// Name is the rule or group name (e.g. "remove-skip", "tape").
	Name string
	// Path is the package import path (used to expand Nested entries).
	Path string
	// Report is func() string — nil for nested plugins.
	Report any
	// Match is func() types.Matcher — nil for traversers and nested plugins.
	Match any
	// Replace is func() types.Replacer — nil for traversers and nested plugins.
	Replace any
	// Traverse is func() types.Traverser — nil for replacers and nested plugins.
	Traverse any
	// Fix is func(ast.Node, []types.Place) — nil for replacers and nested plugins.
	Fix any
	// Rules is types.Nested — non-nil only for nested (grouping) plugins.
	Rules types.Nested
}

// PluginKind is a resolved, runnable plugin.
type PluginKind interface {
	Name() string
	pluginKind()
}

// ReplacerPlugin is a resolved pattern-based plugin.
type ReplacerPlugin struct {
	rule    string
	report  func() string
	match   func() types.Matcher
	replace func() types.Replacer
}

func (p ReplacerPlugin) Name() string            { return p.rule }
func (ReplacerPlugin) pluginKind()               { _ = "replacer" }
func (p ReplacerPlugin) Report() string          { return p.report() }
func (p ReplacerPlugin) Match() types.Matcher    { return p.match() }
func (p ReplacerPlugin) Replace() types.Replacer { return p.replace() }

// TraverserPlugin is a resolved AST-walk plugin.
type TraverserPlugin struct {
	rule     string
	report   func() string
	traverse func() types.Traverser
	fix      func(ast.Node, []types.Place)
}

// RuleState is the enabled/disabled state of a rule from config.
type RuleState struct {
	Enabled bool
	Msg     string
	Options map[string]any
}

// Config maps rule name → state. Parsed from .indra.toml [rules] section.
type Config map[string]RuleState

// DefaultConfig returns a config with all rules enabled (empty filter).
func DefaultConfig() Config { return Config{} }

// candidate is a resolved plugin before config filtering.
type candidate struct {
	rule    string
	kind    PluginKind
	enabled bool // default-enabled from Nested PluginEntry (true for strings)
}

// Load resolves top-level plugins and nested sub-plugins into runnable kinds,
// then filters them by cfg.
//
// A plugin entry whose Rules field is non-nil is a nested group: Load expands
// each of its sub-paths into "group/rule" candidates. Top-level entries keep
// their own rule name.
//
// Filtering priority:
//  1. exact config match: "tape/remove-skip" disabled → rule disabled
//  2. prefix config match: "tape" disabled → all "tape/*" disabled
//  3. PluginEntry.Enabled=false (Off() in Nested) → disabled unless config says on
//  4. default: enabled
func Load(plugins []PluginFuncs, cfg Config) []PluginKind {
	var cands []candidate
	for _, p := range plugins {
		if p.Rules != nil {
			for name, v := range p.Rules {
				path := entryPath(v)
				pf, ok := findPluginFuncs(plugins, path)
				if !ok {
					continue
				}
				rule := p.Name + "/" + name
				cands = append(cands, candidate{rule: rule, kind: resolveKind(pf, rule), enabled: isEntryEnabled(v)})
			}
			continue
		}
		cands = append(cands, candidate{rule: p.Name, kind: resolveKind(p, p.Name), enabled: true})
	}

	out := make([]PluginKind, 0, len(cands))
	for _, c := range cands {
		if !isEnabled(c, cfg) {
			continue
		}
		out = append(out, c.kind)
	}
	return out
}

// isEnabled decides whether a candidate survives config filtering.
func isEnabled(c candidate, cfg Config) bool {
	// 1. exact match
	if st, ok := cfg[c.rule]; ok {
		return st.Enabled
	}
	// 2. prefix match: a disabled group key disables every group/* rule
	for key, st := range cfg {
		if len(key) < len(c.rule) && c.rule[:len(key)] == key && c.rule[len(key)] == '/' && !st.Enabled {
			return false
		}
	}
	// 3/4. default from Nested PluginEntry, else enabled
	return c.enabled
}

func (p TraverserPlugin) Name() string                            { return p.rule }
func (TraverserPlugin) pluginKind()                               { _ = "traverser" }
func (p TraverserPlugin) Report() string                          { return p.report() }
func (p TraverserPlugin) Traverse() types.Traverser               { return p.traverse() }
func (p TraverserPlugin) Fix(node ast.Node, places []types.Place) { p.fix(node, places) }

// resolveKind detects a plugin's kind from its func fields, naming it rule,
// and panics on a malformed shape at init time.
func resolveKind(p PluginFuncs, rule string) PluginKind {
	report := invokeReport(p)
	if p.Match != nil && p.Replace != nil {
		return ReplacerPlugin{
			rule:    rule,
			report:  report,
			match:   mustFunc[func() types.Matcher](p, "Match"),
			replace: mustFunc[func() types.Replacer](p, "Replace"),
		}
	}
	if p.Traverse != nil && p.Fix != nil {
		return TraverserPlugin{
			rule:     rule,
			report:   report,
			traverse: mustFunc[func() types.Traverser](p, "Traverse"),
			fix:      mustFunc[func(ast.Node, []types.Place)](p, "Fix"),
		}
	}
	panic("engine-loader: " + p.Name + ": unknown plugin kind (need Match+Replace or Traverse+Fix)")
}

// invokeReport extracts the func() string Report value.
func invokeReport(p PluginFuncs) func() string {
	if p.Report == nil {
		panic("engine-loader: " + p.Name + ": missing Report func")
	}
	return mustFunc[func() string](p, "Report")
}

// funcValue validates the named field is a non-nil func with an expected shape.
// Report/Match/Replace/Traverse are zero-arg single-return; Fix is two-arg.
func funcValue(p PluginFuncs, field string) reflect.Value {
	raw := reflect.ValueOf(fieldOf(p, field))
	if raw.Kind() != reflect.Func {
		panic("engine-loader: " + p.Name + ": " + field + " is not a func")
	}
	if field == "Fix" {
		if raw.Type().NumIn() != 2 || raw.Type().NumOut() != 0 {
			panic("engine-loader: " + p.Name + ": Fix must be func(ast.Node, []types.Place)")
		}
		return raw
	}
	if raw.Type().NumIn() != 0 || raw.Type().NumOut() != 1 {
		panic("engine-loader: " + p.Name + ": " + field + " must be func() single-return")
	}
	return raw
}

// fieldOf returns the raw value of a named PluginFuncs field.
func fieldOf(p PluginFuncs, field string) any {
	switch field {
	case "Report":
		return p.Report
	case "Match":
		return p.Match
	case "Replace":
		return p.Replace
	case "Traverse":
		return p.Traverse
	case "Fix":
		return p.Fix
	default:
		panic("engine-loader: unknown field " + field)
	}
}

// isEntryEnabled reports the default enabled state of a Nested value.
// A plain string path is enabled; a PluginEntry carries its own state.
func isEntryEnabled(v any) bool {
	if e, ok := v.(types.PluginEntry); ok {
		return e.Enabled
	}
	return true
}

// entryPath extracts the package path from a Nested value (string or PluginEntry).
func entryPath(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case types.PluginEntry:
		return t.Path
	default:
		return ""
	}
}

// findPluginFuncs returns the PluginFuncs whose Path matches, else ok=false.
func findPluginFuncs(plugins []PluginFuncs, path string) (PluginFuncs, bool) {
	for _, pf := range plugins {
		if pf.Path == path {
			return pf, true
		}
	}
	return PluginFuncs{}, false
}

// mustFunc asserts the named field has a func type matching T and returns it.
func mustFunc[T any](p PluginFuncs, field string) T {
	raw := funcValue(p, field)
	fn, ok := raw.Interface().(T)
	if !ok {
		panic("engine-loader: " + p.Name + ": " + field + " has wrong signature")
	}
	return fn
}
