package formatter_test

import (
	"fmt"
	"testing"

	"coderaiser/indra/internal/formatter"
	formatter_codeframe "coderaiser/indra/internal/formatter_codeframe"
	dump "coderaiser/indra/internal/formatter_dump"
	formatter_frame "coderaiser/indra/internal/formatter_frame"
	formatter_json "coderaiser/indra/internal/formatter_json"
	formatter_json_lines "coderaiser/indra/internal/formatter_json_lines"
	formatter_memory "coderaiser/indra/internal/formatter_memory"
	formatter_progress "coderaiser/indra/internal/formatter_progress"
	pb "coderaiser/indra/internal/formatter_progress_bar"
	formatter_stream "coderaiser/indra/internal/formatter_stream"
	formatter_time "coderaiser/indra/internal/formatter_time"

	. "github.com/coderaiser/go-tape"
)

func TestChooseByName(t *testing.T) {
	Test(t, "formatter: ChooseByName json returns json", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("json")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_json.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName time returns time", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("time")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_time.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName memory returns memory", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("memory")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_memory.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName frame returns frame", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("frame")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_frame.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName codeframe returns codeframe", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("codeframe")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_codeframe.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName stream returns stream", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("stream")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_stream.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName progress returns progress", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("progress")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_progress.Format))
		t.End()
	})

	Test(t, "formatter: ChooseByName json-lines returns json-lines", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("json-lines")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_json_lines.Format))
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

	Test(t, "formatter: ChooseByName explicit name beats CI=true", func(t *T) {
		t.TB().Setenv("CI", "true")
		f := formatter.ChooseByName("json-lines")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formatter_json_lines.Format))
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
		t.Equal(result, fmt.Sprintf("%p", formatter_json_lines.Format))

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
