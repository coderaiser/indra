// Package formatter_progress_bar renders lint findings as a live progress bar
// on stderr, falling back to the dump format on the final file.
package formatter_progress_bar

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	dump "coderaiser/indra/internal/formatter_dump"
	"coderaiser/indra/types"
)

const (
	barWidth     = 40
	barComplete  = '█'
	barEmpty     = '░'
	defaultColor = "#6fbdf1"

	// HideCursor and ShowCursor are ANSI escape sequences exported for tests.
	HideCursor = "\033[?25l"
	ShowCursor = "\033[?25h"
)

// Config holds runtime options for the progress bar.
type Config struct {
	Color    string // hex e.g. "#6fbdf1". Empty = defaultColor.
	MinCount int    // show when file count >= MinCount. 0 = always.
}

var cfg = Config{
	Color:    defaultColor,
	MinCount: 0,
}

// Configure sets the active progress bar configuration.
// Call once at program startup before any Format call.
func Configure(c Config) {
	if c.Color != "" {
		cfg.Color = c.Color
	}
	cfg.MinCount = c.MinCount
}

// Format is the progress-bar formatter. It writes a live bar to stderr
// mid-run and returns the dump output on the last file.
func Format(name string, _ []byte, places []types.Place, index, count, filesWithIssues, errorsCount int) string {
	result := dump.Format(name, nil, places, index, count, filesWithIssues, errorsCount)

	if !ShouldShow(count) {
		return result
	}

	if index == 0 {
		fmt.Fprintf(os.Stderr, "%s", HideCursor)
	}

	errStr := "👌"
	if errorsCount > 0 {
		errStr = fmt.Sprintf("\033[31m%d\033[0m", errorsCount)
	}

	bar := RenderBar(index+1, count, cfg.Color)
	pct := 0
	if count > 0 {
		pct = (index + 1) * 100 / count
	}
	line := fmt.Sprintf("%s %d%% | %s | %d/%d | %s",
		bar, pct, errStr, index+1, count, Truncate(name, 40))
	width := TermWidth()
	if VisibleLen(line) > width {
		line = TruncateANSI(line, width)
	}
	fmt.Fprintf(os.Stderr, "\r%s", line)

	if index == count-1 {
		fmt.Fprintf(os.Stderr, "\r\033[2K%s", ShowCursor)
		return result
	}
	return ""
}

// ShouldShow returns true if the progress bar should be displayed.
func ShouldShow(count int) bool {
	if os.Getenv("INDRA_PROGRESS_BAR") == "1" {
		return true
	}
	if os.Getenv("INDRA_PROGRESS_BAR") == "0" {
		return false
	}
	if ci := os.Getenv("CI"); ci == "1" || ci == "true" {
		return false
	}
	return count >= cfg.MinCount
}

// RenderBar renders a Unicode block progress bar. Exported for testing.
func RenderBar(done, total int, color string) string {
	ansi := hexToANSI(color)
	if total == 0 {
		return fmt.Sprintf("%s%s\033[0m", ansi, strings.Repeat(string(barEmpty), barWidth))
	}
	filled := done * barWidth / total
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat(string(barComplete), filled) +
		strings.Repeat(string(barEmpty), barWidth-filled)
	return fmt.Sprintf("%s%s\033[0m", ansi, bar)
}

func hexToANSI(color string) string {
	if len(color) != 7 || color[0] != '#' {
		return color
	}
	parse := func(s string) int {
		n := 0
		for _, c := range s {
			n <<= 4
			switch {
			case c >= '0' && c <= '9':
				n |= int(c - '0')
			case c >= 'a' && c <= 'f':
				n |= int(c-'a') + 10
			case c >= 'A' && c <= 'F':
				n |= int(c-'A') + 10
			}
		}
		return n
	}
	r, g, b := parse(color[1:3]), parse(color[3:5]), parse(color[5:7])
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

// VisibleLen returns the visible character length excluding ANSI escapes.
func VisibleLen(s string) int {
	return len([]rune(ansiEscape.ReplaceAllString(s, "")))
}

// TruncateANSI truncates s to n visible characters, preserving ANSI codes.
func TruncateANSI(s string, n int) string {
	plain := ansiEscape.ReplaceAllString(s, "")
	if len([]rune(plain)) <= n {
		return s
	}
	return string([]rune(plain)[:n])
}

// Truncate truncates s to n runes, adding "..." if truncated.
func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-3]) + "..."
}
