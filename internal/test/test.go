// Package test shims the exported acceptance-test helpers to indra's engine.
// Plugin test files keep their existing dot-import, so moving the engine
// contract behind types.Lint changes nothing for them.
package test

import (
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

// ForGroup is the exported ForGroup bound to indra's engine.
var ForGroup = indratest.ForGroup(indraLint)

// indraLint converts []any plugin args into runnable PluginItems and runs
// indra's engine over src.
func indraLint(src []byte, fix bool, plugins []any) (types.LintResult, error) {
	items := make([]runner.PluginItem, 0, len(plugins))
	for _, payload := range plugins {
		arg := payload.(indratest.PluginArg)
		kinds := loader.Load([]loader.PluginFuncs{{Name: arg.Rule, Plugin: arg.Plugin}}, loader.Config{})
		validatePlugin(kinds[0])
		items = append(items, runner.PluginItem{Rule: kinds[0].Name(), Plugin: kinds[0]})
	}
	result, err := processor.Process(processor.Params{Src: src, Fix: fix, Plugins: items})
	return types.LintResult(result), err
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
