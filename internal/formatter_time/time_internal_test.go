package formatter_time

import (
	"testing"
	"time"

	. "github.com/coderaiser/go-tape"
)

func TestFormatDuration(t *testing.T) {
	Test(t, "time: formatDuration seconds", func(t *T) {
		result := formatDuration(1500 * time.Millisecond)
		t.Equal(result, "1.50s")

		t.End()
	})

	Test(t, "time: formatDuration milliseconds", func(t *T) {
		result := formatDuration(5 * time.Millisecond)
		t.Equal(result, "5ms")

		t.End()
	})

	Test(t, "time: formatDuration microseconds", func(t *T) {
		result := formatDuration(10 * time.Microsecond)
		t.Equal(result, "10µs")

		t.End()
	})
}
