//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

// extract-verb-skips: the callback body has non-matching statements before
// the first t.Report call — a non-ExprStmt (x := 1), an ExprStmt that is not
// a CallExpr (y), and a SelectorExpr call whose receiver is not "t" (s.Do()).
// extractVerb skips each of these with the corresponding continue branch.
var Test = CreateTest("extract-verb-skips", nil)

func f(t *testing.T) {
	Test(t, "extract-verb-skips: report: extract-verb-skips", func(t *T) {
		s.Do()
		y
		x := 1
		t.Report("extract-verb-skips")
	})
}
