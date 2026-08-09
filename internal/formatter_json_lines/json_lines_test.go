package formatter_json_lines_test

import (
	"encoding/json"
	"strings"
	"testing"

	formatter_json_lines "coderaiser/indra/internal/formatter_json_lines"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var places1 = []types.Place{
	{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 10, Column: 3}},
}

func TestFormatterJsonLines(t *testing.T) {
	Test(t, "json-lines: output is newline terminated", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.HasSuffix(out, "\n"))
		t.End()
	})

	Test(t, "json-lines: output is valid JSON", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 0, 1, 0, 0)
		var m map[string]any
		err := json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.NotOk(err)

		t.End()
	})

	Test(t, "json-lines: name field is correct", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 0, 1, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["name"], "foo.go")
		t.End()
	})

	Test(t, "json-lines: empty places field for clean file", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 0, 1, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		result := len(m["places"].([]any))
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "json-lines: places field contains findings", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, places1, 0, 1, 1, 1)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		result := len(m["places"].([]any))
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "json-lines: index field is correct", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 2, 5, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["index"], float64(2))
		t.End()
	})

	Test(t, "json-lines: count field is correct", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, nil, 2, 5, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["count"], float64(5))
		t.End()
	})

	Test(t, "json-lines: filesCount field is correct", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, places1, 0, 3, 1, 2)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["filesCount"], float64(1))
		t.End()
	})

	Test(t, "json-lines: errorsCount field is correct", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, places1, 0, 3, 1, 2)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["errorsCount"], float64(2))
		t.End()
	})

	Test(t, "json-lines: one record per file call", func(t *T) {
		out1 := formatter_json_lines.Format("a.go", nil, nil, 0, 2, 0, 0)
		out2 := formatter_json_lines.Format("b.go", nil, places1, 1, 2, 1, 1)
		lines := strings.Split(strings.TrimSpace(out1+out2), "\n")
		result := len(lines)
		t.Equal(result, 2)

		t.End()
	})

	Test(t, "json-lines: place rule key is lowercase (putout-compatible)", func(t *T) {
		out := formatter_json_lines.Format("foo.go", nil, places1, 0, 1, 1, 1)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		first := m["places"].([]any)[0].(map[string]any)
		t.Equal(first["rule"], "tape/remove-skip")
		t.End()
	})
}
