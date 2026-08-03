package remove_useless_match

import . "coderaiser/indra/types"

func Report() string { return "remove useless Match" }

func Replace() Replacer {
	return Replacer{
		// Matcher with all-nil guards — guard is a no-op, Match can be deleted
		`func Match() Matcher { return Matcher{__a: nil} }`: "",
		// Empty Matcher — no patterns, Match does nothing
		`func Match() Matcher { return Matcher{} }`: "",
	}
}
