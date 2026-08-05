package formatter_test

import (
	"fmt"
	"testing"

	"coderaiser/indra/internal/formatter"
	dump "coderaiser/indra/internal/formatter_dump"
	jsonlines "coderaiser/indra/internal/formatter_json_lines"
	pb "coderaiser/indra/internal/formatter_progress_bar"

	. "github.com/coderaiser/go-tape"
)

func TestChoose(t *testing.T) {
	Test(t, "formatter: CI=true returns dump", func(t *T) {
		t.TB().Setenv("CI", "true")
		t.TB().Setenv("INDRA_FORMATTER", "")
		f := formatter.Choose()
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", dump.Format))

		t.End()
	})

	Test(t, "formatter: INDRA_FORMATTER=dump returns dump", func(t *T) {
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_FORMATTER", "dump")
		f := formatter.Choose()
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", dump.Format))

		t.End()
	})

	Test(t, "formatter: INDRA_FORMATTER=json-lines returns json-lines", func(t *T) {
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_FORMATTER", "json-lines")
		f := formatter.Choose()
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", jsonlines.Format))

		t.End()
	})

	Test(t, "formatter: default returns progress-bar", func(t *T) {
		t.TB().Setenv("CI", "")
		t.TB().Setenv("INDRA_FORMATTER", "")
		f := formatter.Choose()
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", pb.Format))

		t.End()
	})
}
