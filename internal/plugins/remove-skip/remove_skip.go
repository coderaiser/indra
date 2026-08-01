package removeskip

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:   "remove-skip",
	Report: func() string { return "remove t.Skip call" },
	Match: func() map[string]engine.MatchFn {
		return map[string]engine.MatchFn{
			"__recv.Skip(__args)": nil,
		}
	},
}
