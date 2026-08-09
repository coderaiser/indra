// Package formatter selects the active output formatter based on the
// environment.
package formatter

import (
	"os"

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
		return formatter_json.Format
	case "json-lines":
		return formatter_json_lines.Format
	case "progress":
		return formatter_progress.Format
	case "codeframe":
		return formatter_codeframe.Format
	case "frame":
		return formatter_frame.Format
	case "memory":
		return formatter_memory.Format
	case "time":
		return formatter_time.Format
	case "stream":
		return formatter_stream.Format
	case "dump":
		return dump.Format
	default:
		return pb.Format
	}
}
