package formatter_test

import (
	"fmt"
	"testing"

	"coderaiser/indra/internal/formatter"
	formcf "coderaiser/indra/internal/formatter_codeframe"
	dump "coderaiser/indra/internal/formatter_dump"
	formframe "coderaiser/indra/internal/formatter_frame"
	formjson "coderaiser/indra/internal/formatter_json"
	jsonlines "coderaiser/indra/internal/formatter_json_lines"
	formmem "coderaiser/indra/internal/formatter_memory"
	formprog "coderaiser/indra/internal/formatter_progress"
	pb "coderaiser/indra/internal/formatter_progress_bar"
	formstream "coderaiser/indra/internal/formatter_stream"
	formtime "coderaiser/indra/internal/formatter_time"

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

	Test(t, "formatter: ChooseByName time returns time", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("time")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formtime.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName memory returns memory", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("memory")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formmem.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName frame returns frame", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("frame")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formframe.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName codeframe returns codeframe", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("codeframe")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formcf.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName stream returns stream", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("stream")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formstream.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName progress returns progress", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("progress")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formprog.Format))
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
