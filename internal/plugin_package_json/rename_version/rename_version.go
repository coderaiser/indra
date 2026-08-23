package rename_version

import . "coderaiser/indra/types"

func Report() string { return "normalise version field" }

func Replace() Replacer {
	return Replacer{
		`"version": "__a"`: `"version": "0.0.0"`,
	}
}

type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
