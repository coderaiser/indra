package formatter_time

import (
	"testing"
	"time"

	. "github.com/coderaiser/go-tape"
)

func TestFormatDuration(t *testing.T) {
	Test(t, "time: formatDuration seconds", func(t *T) {
		t.Equal(formatDuration(1500*time.Millisecond), "1.50s")
		t.End()
	})

	Test(t, "time: formatDuration milliseconds", func(t *T) {
		t.Equal(formatDuration(5*time.Millisecond), "5ms")
		t.End()
	})

	Test(t, "time: formatDuration microseconds", func(t *T) {
		t.Equal(formatDuration(10*time.Microsecond), "10µs")
		t.End()
	})
}
