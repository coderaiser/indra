package formatter_test

import (
	"fmt"
	"testing"

	"coderaiser/indra/internal/formatter"
	dump "coderaiser/indra/internal/formatter_dump"
	formjson "coderaiser/indra/internal/formatter_json"
	jsonlines "coderaiser/indra/internal/formatter_json_lines"
	pb "coderaiser/indra/internal/formatter_progress_bar"

	. "github.com/coderaiser/go-tape"
)

func TestChooseByName(t *testing.T) {
	Test(t, "formatter: ChooseByName json returns json", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("json")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formjson.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName json-lines returns json-lines", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("json-lines")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", jsonlines.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName dump returns dump", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("dump")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", dump.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName unknown returns progress-bar", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("unknown")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", pb.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName CI=true returns dump regardless of name", func(t *T) {
		t.TB().Setenv("CI", "true")
		f := formatter.ChooseByName("json-lines")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", dump.Format))
		t.End()
	})
}

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
