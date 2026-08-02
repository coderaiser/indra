package removeskip

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:    "remove-skip",
	Report:  report,
	Match:   match,
	Replace: replace,
}

func report() string { return "remove Test.Skip call" }

func match() map[string]engine.MatchFn {
	return map[string]engine.MatchFn{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: nil,
	}
}

func replace() map[string]string {
	return map[string]string{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
