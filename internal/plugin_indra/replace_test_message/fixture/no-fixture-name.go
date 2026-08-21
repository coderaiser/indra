//go:build ignore

package fixture

import (
	"testing"

	. "coderaiser/indra/internal/test"
)

// no-fixture-name: the second Test call has a report verb but its t.Report
// carries no string argument, so no fixture name is derivable and fixTestCall
// leaves that call untouched while the first call still gets fixed.
var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: report: wrong", func(t *T) {
		t.Report("some-fix", "remove Test.Skip call")
		t.End()
	})
	Test(t, "remove-skip: report: manual", func(t *T) {
		t.Report()
		t.End()
	})
}
