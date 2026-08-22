// Package engine_loader resolves plugin packages into runnable plugin kinds
// and filters them by rule configuration.
package engine_loader

import (
	"reflect"

	"coderaiser/indra/types"
)

// PluginFuncs is one entry of the plugin registry. It never describes a
// plugin's shape (Report/Match/Replace/Traverse/Fix): a leaf carries the whole
// method-bearing struct in Plugin, and a group carries its []types.Rule. The
// loader detects the shape through reflection, so a plugin can switch between
// replacer and traverser without touching the registry.
type PluginFuncs struct {
	// Name is the rule or group name (e.g. "remove-skip", "tape").
	Name string
	// Plugin is a leaf plugin — a struct exposing Report/Match/Replace/Traverse/
	// Fix methods. Nil for groups.
	Plugin any
	// Rules is a group's sub-rules. Non-nil means this is a nested group.
	Rules []types.Rule
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
	report   types.ReportFn
	traverse func() types.Traverser
	fix      types.FixFn
}

func (p TraverserPlugin) Name() string                              { return p.rule }
func (TraverserPlugin) pluginKind()                                 { _ = "traverser" }
func (p TraverserPlugin) Report(pPath types.Path) string            { return p.report(pPath) }
func (p TraverserPlugin) Traverse() types.Traverser                 { return p.traverse() }
func (p TraverserPlugin) Fix(pPath types.Path, opts map[string]any) { p.fix(pPath, opts) }

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
	rule string
	kind PluginKind
}

// Load resolves top-level plugins and nested sub-plugins into runnable kinds,
// then filters them by cfg.
//
// A plugin entry whose Rules field is non-nil is a nested group: Load expands
// each rule into "group/rule" candidates. Top-level entries keep their own name.
//
// Filtering priority:
//  1. exact config match: "tape/remove-skip" disabled → rule disabled
//  2. prefix config match: "tape" disabled → all "tape/*" disabled
//  3. default: enabled
func Load(plugins []PluginFuncs, cfg Config) []PluginKind {
	var cands []candidate
	for _, p := range plugins {
		if p.Rules != nil {
			for _, r := range p.Rules {
				rule := p.Name + "/" + r.Name
				cands = append(cands, candidate{rule: rule, kind: resolve(r.Plugin, rule, p.Name, ruleOptions(cfg, rule))})
			}
			continue
		}
		cands = append(cands, candidate{rule: p.Name, kind: resolve(p.Plugin, p.Name, p.Name, ruleOptions(cfg, p.Name))})
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
	// 2. group prefix match: "tape" off → all "tape/*" disabled; "tape" on → all
	// "tape/*" enabled.
	for key, st := range cfg {
		if len(key) < len(c.rule) && c.rule[:len(key)] == key && c.rule[len(key)] == '/' {
			return st.Enabled
		}
	}
	// 3. default: enabled
	return true
}

// resolve detects a plugin's kind from its exported methods, naming it rule.
// A plugin with a Replace method is a replacer; one with a Traverse method is
// a traverser. It panics on a malformed shape at init time.
// A replacer may expose Filter() types.Filter, Match(Options) types.Matcher,
// or Match() types.Matcher. Priority: a non-empty Filter wins (includer
// pattern); otherwise Match(Options) (putout-aligned, closes over the rule's
// Options); otherwise a plain Match(); otherwise an empty matcher.
func resolve(plugin any, rule, owner string, opts types.Options) PluginKind {
	v := reflect.ValueOf(plugin)
	if r, ok := method[func() types.Replacer](v, "Replace"); ok {
		match := matchFactory(v, opts)
		if f, ok := method[func() types.Filter](v, "Filter"); ok && len(f()) > 0 {
			match = filterToMatch(f(), opts)
		}
		return ReplacerPlugin{
			rule:    rule,
			report:  mustMethod[func() string](v, "Report"),
			match:   match,
			replace: r,
		}
	}
	if tr, ok := method[func() types.Traverser](v, "Traverse"); ok {
		return TraverserPlugin{
			rule:     rule,
			report:   mustMethod[types.ReportFn](v, "Report"),
			traverse: tr,
			fix:      mustMethod[types.FixFn](v, "Fix"),
		}
	}
	panic("engine-loader: " + owner + ": unknown plugin kind (need Replace or Traverse)")
}

// matchFactory returns a Matcher factory for the plugin's Match method.
// Priority: Match(Options) types.Matcher closes over opts (putout-aligned);
// otherwise the plain Match() types.Matcher; otherwise an empty matcher.
// Signature detection is probed by reflect type rather than the strict
// `method` helper, so a plugin with a plain Match() (0-arg) is not treated as
// a malformed Match(Options) (1-arg). A Match that is neither signature — e.g.
// Match() int — still panics so a broken plugin is caught at load time.
func matchFactory(v reflect.Value, opts types.Options) func() types.Matcher {
	m := v.MethodByName("Match")
	if !m.IsValid() {
		return func() types.Matcher { return types.Matcher{} }
	}
	switch m.Type() {
	case reflect.TypeOf(types.MatchWithOpts(nil)):
		return func() types.Matcher { return m.Interface().(types.MatchWithOpts)(opts) }
	case reflect.TypeOf(func() types.Matcher(nil)):
		return func() types.Matcher { return m.Interface().(func() types.Matcher)() }
	}
	panic("engine-loader: method Match has wrong signature")
}

// filterToMatch wraps a Filter into a Matcher factory, closing over the rule's
// Options so each FilterFn receives them at guard time.
func filterToMatch(f types.Filter, opts types.Options) func() types.Matcher {
	return func() types.Matcher {
		m := make(types.Matcher, len(f))
		for pattern, fn := range f {
			m[pattern] = func(vars types.Vars, p types.Path) bool {
				return fn(vars, p, opts)
			}
		}
		return m
	}
}

// ruleOptions returns the configured Options for rule, or nil.
func ruleOptions(cfg Config, rule string) types.Options {
	if st, ok := cfg[rule]; ok {
		return st.Options
	}
	return nil
}

// method returns the named method of plugin value v asserted to type T.
// It reports false when the method is absent; it panics when the method exists
// with an incompatible signature.
func method[T any](v reflect.Value, name string) (T, bool) {
	var zero T
	if !v.IsValid() {
		return zero, false
	}
	m := v.MethodByName(name)
	if !m.IsValid() {
		return zero, false
	}
	fn, ok := m.Interface().(T)
	if !ok {
		panic("engine-loader: method " + name + " has wrong signature")
	}
	return fn, true
}

// mustMethod is method but panics when the method is missing.
func mustMethod[T any](v reflect.Value, name string) T {
	fn, ok := method[T](v, name)
	if !ok {
		panic("engine-loader: missing " + name + " method")
	}
	return fn
}
