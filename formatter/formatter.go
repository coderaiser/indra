// Package formatter selects the active output formatter based on the
// environment. External linters import only this package.
package formatter

import internal "coderaiser/indra/internal/formatter"

// Func is called once per file with running totals. See internal/formatter
// for the full contract.
type Func = internal.Func

// Choose returns the active formatter based on environment.
// CI=true → dump. INDRA_FORMATTER=<name> → that formatter. Default: progress-bar.
func Choose() Func { return internal.Choose() }

// ChooseByName returns the formatter for the given name.
// Valid names: json, json-lines, progress, codeframe, frame, memory, time,
// stream, dump. Empty string or unknown falls back to Choose().
func ChooseByName(name string) Func { return internal.ChooseByName(name) }
