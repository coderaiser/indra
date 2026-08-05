package remove_useless_condition

import . "coderaiser/indra/types"

func Report() string { return "remove useless condition: Ok(err != nil) → Ok(err)" }

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b != nil)": "__a.Ok(__b)",
	}
}
