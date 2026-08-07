// Package formatter_json accumulates all findings and emits a single
// pretty-printed JSON blob on the last file.
package formatter_json

import (
	"encoding/json"
	"sync"

	"coderaiser/indra/types"
)

type fileError struct {
	Name   string        `json:"name"`
	Places []types.Place `json:"places"`
}

type result struct {
	Errors      []fileError `json:"errors"`
	FilesCount  int         `json:"filesCount"`
	ErrorsCount int         `json:"errorsCount"`
}

var (
	mu     sync.Mutex
	errors []fileError
)

// Format accumulates files with places and returns one JSON blob on the last
// file, resetting its internal state afterwards.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	mu.Lock()
	defer mu.Unlock()

	if len(places) > 0 {
		errors = append(errors, fileError{Name: name, Places: places})
	}

	if index != count-1 {
		return ""
	}

	r := result{
		Errors:      errors,
		FilesCount:  filesWithIssues,
		ErrorsCount: errorsCount,
	}
	if r.Errors == nil {
		r.Errors = []fileError{}
	}
	errors = nil // reset

	b, _ := json.MarshalIndent(r, "", "    ")
	return string(b) + "\n"
}
