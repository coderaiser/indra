// Package test shims the exported acceptance-test helpers to indra's engine.
// Plugin test files keep their existing dot-import, so moving the engine
// contract behind types.Lint changes nothing for them.
package test

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"

	loader "coderaiser/indra/engine_loader"
	processor "coderaiser/indra/engine_processor"
	runner "coderaiser/indra/engine_runner"
	indratest "coderaiser/indra/test"
	"coderaiser/indra/types"
)

// T is the test runner type from the exported package, bound to indra.
type T = indratest.T

// CreateTest is the exported CreateTest bound to indra's engine.
var CreateTest = indratest.CreateTest(indraLint)

// indraLint converts []any plugin args into runnable PluginItems and runs
// indra's engine over src.
func indraLint(src []byte, fix bool, plugins []any) (types.LintResult, error) {
	items := make([]runner.PluginItem, 0, len(plugins))
	for _, payload := range plugins {
		arg := payload.(indratest.PluginArg)
		cfg := arg.Config
		if cfg == nil {
			cfg = loader.Config{}
		}
		kinds := loader.Load([]loader.PluginFuncs{{Name: arg.Rule, Plugin: arg.Plugin}}, cfg)
		validatePlugin(kinds[0])
		items = append(items, runner.PluginItem{Rule: kinds[0].Name(), Plugin: kinds[0]})
	}
	result, err := processor.Process(processor.Params{Src: src, Fix: fix, Plugins: items})
	if err != nil {
		// Non-Go sources (e.g. package.json fixtures) have no Go AST —
		// fall back to plain-text replacer application. Without any
		// replacer plugin there is nothing text mode can do: surface the
		// parse error.
		for _, item := range items {
			if _, ok := item.Plugin.(loader.ReplacerPlugin); ok {
				return textFallback(src, items, fix)
			}
		}
	}
	return types.LintResult(result), err
}

// textFallback applies replacer patterns to src as plain text. Placeholder
// tokens (__a, __b, …) in a pattern match any, possibly empty, run of
// characters; the matched text is substituted into the template.
func textFallback(src []byte, items []runner.PluginItem, fix bool) (types.LintResult, error) {
	out := string(src)
	var places []types.Place
	for _, item := range items {
		replacer, ok := item.Plugin.(loader.ReplacerPlugin)
		if !ok {
			continue
		}
		for pattern, tmpl := range replacer.Replace() {
			loc, groups := findTextPattern(out, pattern)
			if loc == nil {
				continue
			}
			matchText := out[loc[0]:loc[1]]
			replacement := expandTemplate(groups, tmpl)
			if replacement == matchText {
				// The rewrite is an identity — nothing to report or fix.
				continue
			}
			places = append(places, types.Place{
				Rule:     item.Rule,
				Message:  replacer.Report(),
				Position: lineColumn(out, loc[0]),
			})
			if fix {
				out = out[:loc[0]] + replacement + out[loc[1]:]
			}
		}
	}
	result := types.LintResult{Places: places}
	if fix {
		result.Out = []byte(out)
	} else {
		result.Out = src
	}
	return result, nil
}

// placeholderRe matches __a-style wildcard tokens in patterns and templates.
var placeholderRe = regexp.MustCompile(`__[a-z]`)

// findTextPattern locates the first occurrence of pattern in src, returning
// the match location and the captured placeholder texts. Returns a nil
// location when the pattern does not occur.
func findTextPattern(src string, pattern string) ([]int, []string) {
	parts := placeholderRe.Split(pattern, -1)
	count := len(parts) - 1
	var sb strings.Builder
	for i, part := range parts {
		sb.WriteString(regexp.QuoteMeta(part))
		if i < count {
			sb.WriteString(`(?s:(.*?))`)
		}
	}
	re, err := regexp.Compile(sb.String())
	if err != nil {
		return nil, nil
	}
	loc := re.FindStringSubmatchIndex(src)
	if loc == nil {
		return nil, nil
	}
	var groups []string
	for i := 0; i < count; i++ {
		groups = append(groups, src[loc[2*i+2]:loc[2*i+3]])
	}
	return loc[:2], groups
}

// expandTemplate substitutes the captured texts into tmpl, replacing each
// __a-style token with its capture in order.
func expandTemplate(groups []string, tmpl string) string {
	i := 0
	return placeholderRe.ReplaceAllStringFunc(tmpl, func(string) string {
		if i >= len(groups) {
			return ""
		}
		replacement := groups[i]
		i++
		return replacement
	})
}

// lineColumn converts a byte offset into a one-based line and column.
func lineColumn(src string, offset int) types.Position {
	line := 1 + strings.Count(src[:offset], "\n")
	column := offset
	if idx := strings.LastIndexByte(src[:offset], '\n'); idx != -1 {
		column = offset - idx - 1
	}
	return types.Position{Line: line, Column: column + 1}
}

// CreateTestConfig is CreateTest with a loader.Config applied when resolving
// the plugin — used to exercise option-driven rules (Filter + Options).
func CreateTestConfig(rule string, plugin any, cfg loader.Config) func(*testing.T, string, func(*T)) {
	plugins := []any{indratest.PluginArg{Rule: rule, Plugin: plugin, Config: cfg}}
	dir := callerFixtureDir(1)
	return tape.Extend(func(base *tape.T) *T {
		return indratest.New(base, indraLint, plugins, dir)
	})
}

// callerFixtureDir returns the fixture/ directory next to the test file that
// called CreateTestConfig. depth is the position of the caller's frame.
func callerFixtureDir(depth int) string {
	_, file, _, _ := runtime.Caller(depth + 1)
	return filepath.Join(filepath.Dir(file), "fixture")
}

// validatePlugin enforces consistency on a resolved ReplacerPlugin before it is
// used in tests: every Match entry must have a non-nil guard function and every
// Match key must also appear as a Replace key. A malformed plugin would
// otherwise pass the tester silently and fail only at fix time.
func validatePlugin(kind loader.PluginKind) {
	replacer, ok := kind.(loader.ReplacerPlugin)
	if !ok {
		return
	}
	matcher := replacer.Match()
	replacements := replacer.Replace()
	for pattern, guard := range matcher {
		if guard == nil {
			panic("internal/test: " + kind.Name() + ": nil MatchFn for pattern " + pattern)
		}
		if _, ok := replacements[pattern]; !ok {
			panic("internal/test: " + kind.Name() + ": Match key not in Replace: " + pattern)
		}
	}
}
