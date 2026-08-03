package convert_ok_to_not_ok

import . "coderaiser/indra/types"

func Report() string { return "convert Ok(err == nil) to NotOk(err)" }

func Match() Matcher {
	return Matcher{
		"__a.Ok(__b == nil)": nil,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b == nil)": "__a.NotOk(__b)",
	}
}
