// Package formatter selects the active output formatter based on the
// environment.
package formatter

import (
	"os"

	dump "coderaiser/indra/internal/formatter_dump"
	formjson "coderaiser/indra/internal/formatter_json"
	jsonlines "coderaiser/indra/internal/formatter_json_lines"
	formprog "coderaiser/indra/internal/formatter_progress"
	pb "coderaiser/indra/internal/formatter_progress_bar"
	formstream "coderaiser/indra/internal/formatter_stream"
	"coderaiser/indra/types"
)

// Func is called once per file with running totals.
// source is the file's raw content (nil if unread).
// index is 0-based. count is total number of files.
// filesWithIssues and errorsCount are running totals before this call.
// Detects last file via index == count-1.
// Returns the string to write to output (empty = nothing to write yet).
type Func func(name string, source []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string

// Choose returns the active formatter based on environment.
// CI=true → dump. INDRA_FORMATTER=json-lines|dump → that one. Default: progress-bar.
func Choose() Func {
	return ChooseByName(os.Getenv("INDRA_FORMATTER"))
}

// ChooseByName returns the formatter for the given name.
// Empty string or unknown name falls back to env/CI logic.
func ChooseByName(name string) Func {
	if os.Getenv("CI") == "true" {
		return dump.Format
	}
	switch name {
	case "json":
		return formjson.Format
	case "json-lines":
		return jsonlines.Format
	case "progress":
		return formprog.Format
	case "stream":
		return formstream.Format
	case "dump":
		return dump.Format
	default:
		return pb.Format
	}
}
