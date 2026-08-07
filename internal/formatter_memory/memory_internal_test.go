package formatter_memory

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestFormatBytes(t *testing.T) {
	Test(t, "memory: formatBytes small returns bytes", func(t *T) {
		result := formatBytes(500)
		t.Equal(result, "500 B")

		t.End()
	})

	Test(t, "memory: formatBytes kilobytes", func(t *T) {
		result := formatBytes(2048)
		t.Equal(result, "2.0 KB")

		t.End()
	})

	Test(t, "memory: formatBytes megabytes", func(t *T) {
		result := formatBytes(3 * 1 << 20)
		t.Equal(result, "3.0 MB")

		t.End()
	})
}
