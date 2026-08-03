package remove_unused_variable

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestReportMessage(t *testing.T) {
	tape.Test(t, "report: returns remove unused variable", func(t *tape.T) {
		t.Equal(Report(), "remove unused variable")
		t.End()
	})
}
