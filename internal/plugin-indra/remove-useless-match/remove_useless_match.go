package remove_useless_match

import . "coderaiser/indra/types"

func Report() string { return "remove useless Match" }

// Match returns an empty Matcher — this plugin operates at declaration level
// via CompareDecl in walkDecls. No stmt-level patterns needed.
func Match() Matcher { return Matcher{} }

func Replace() Replacer {
	return Replacer{
		// Matcher with all-nil guards — guard is a no-op, Match can be deleted
		`func Match() Matcher { return Matcher{__a: nil} }`: "",
		// Empty Matcher — no patterns, Match does nothing
		`func Match() Matcher { return Matcher{} }`:          "",
	}
}
