package config

import (
	"strings"
	"testing"

	. "github.com/coderaiser/go-tape"
)

// TestDefaultPanic drives the malformed-default path by swapping the embedded
// default.toml bytes, then restores them for the rest of the suite.
func TestDefaultPanic(t *testing.T) {
	original := defaultToml
	t.Cleanup(func() { defaultToml = original })
	defaultToml = []byte("[bad")

	Test(t, "config: Default panics on malformed default.toml", func(t *T) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					msg := r.(string)
					t.Ok(strings.Contains(msg, "malformed default.toml"))
				}
			}()
			Default()
		}()

		t.End()
	})
}
