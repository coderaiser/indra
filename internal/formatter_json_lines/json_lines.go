// Package formatter_json_lines renders one JSON record per file as newline
// separated output, suitable for machine consumption.
package formatter_json_lines

import (
	"encoding/json"

	"coderaiser/indra/types"
)

type record struct {
	Name        string        `json:"name"`
	Places      []types.Place `json:"places"`
	Index       int           `json:"index"`
	Count       int           `json:"count"`
	FilesCount  int           `json:"filesCount"`
	ErrorsCount int           `json:"errorsCount"`
}

// Format emits one JSON record per file call, newline-terminated.
// Stateless — no accumulation between calls.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	if places == nil {
		places = []types.Place{}
	}
	b, _ := json.Marshal(record{
		Name:        name,
		Places:      places,
		Index:       index,
		Count:       count,
		FilesCount:  filesWithIssues,
		ErrorsCount: errorsCount,
	})
	return string(b) + "\n"
}
