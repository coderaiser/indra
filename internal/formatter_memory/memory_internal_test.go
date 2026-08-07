package formatter_memory

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestFormatBytes(t *testing.T) {
	Test(t, "memory: formatBytes small returns bytes", func(t *T) {
		t.Equal(formatBytes(500), "500 B")
		t.End()
	})

	Test(t, "memory: formatBytes kilobytes", func(t *T) {
		t.Equal(formatBytes(2048), "2.0 KB")
		t.End()
	})

	Test(t, "memory: formatBytes megabytes", func(t *T) {
		t.Equal(formatBytes(3*1<<20), "3.0 MB")
		t.End()
	})
}
