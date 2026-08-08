package formatter_json_test

import (
	"encoding/json"
	"strings"
	"testing"

	formjson "coderaiser/indra/internal/formatter_json"
	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

var place1 = types.Place{Rule: "tape/remove-skip", Message: "remove Test.Skip call", Position: types.Position{Line: 5, Column: 2}}

func TestJson(t *testing.T) {
	Test(t, "json: returns empty mid-run", func(t *T) {
		out := formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		t.Equal(out, "")
		t.End()
	})

	Test(t, "json: returns JSON on last file", func(t *T) {
		formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		out := formjson.Format("b.go", nil, nil, 1, 2, 1, 1)
		t.Ok(strings.Contains(out, `"errors"`))
		t.End()
	})

	Test(t, "json: errors array contains files with places", func(t *T) {
		formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		out := formjson.Format("b.go", nil, nil, 1, 2, 1, 1)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		errors := m["errors"].([]any)
		result := len(errors)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "json: clean files not in errors array", func(t *T) {
		out := formjson.Format("a.go", nil, nil, 0, 1, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		errors := m["errors"].([]any)
		result := len(errors)
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "json: filesCount in output", func(t *T) {
		formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		out := formjson.Format("b.go", nil, nil, 1, 2, 1, 1)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["filesCount"], float64(1))
		t.End()
	})

	Test(t, "json: errorsCount in output", func(t *T) {
		formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		out := formjson.Format("b.go", nil, nil, 1, 2, 1, 1)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		t.Equal(m["errorsCount"], float64(1))
		t.End()
	})

	Test(t, "json: output is newline terminated", func(t *T) {
		out := formjson.Format("a.go", nil, nil, 0, 1, 0, 0)
		t.Ok(strings.HasSuffix(out, "\n"))
		t.End()
	})

	Test(t, "json: state resets between runs", func(t *T) {
		// run 1: full multi-file run with one error, resets state after output
		formjson.Format("a.go", nil, []types.Place{place1}, 0, 2, 1, 1)
		formjson.Format("b.go", nil, nil, 1, 2, 1, 1)
		// run 2: fresh single clean file must not inherit the prior error
		out := formjson.Format("c.go", nil, nil, 0, 1, 0, 0)
		var m map[string]any
		json.Unmarshal([]byte(strings.TrimSpace(out)), &m)
		errors := m["errors"].([]any)
		result := len(errors)
		t.Equal(result, 0)

		t.End()
	})
}
