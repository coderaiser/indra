package formatter_test

import (
	"testing"

	"coderaiser/indra/formatter"

	. "github.com/coderaiser/go-tape"
)

func TestChoose(t *testing.T) {
	Test(t, "formatter: Choose returns a non-nil Func", func(t *T) {
		f := formatter.Choose()
		t.Ok(f)

		t.End()
	})
}

func TestChooseByName(t *testing.T) {
	names := []string{"json", "json-lines", "progress", "codeframe",
		"frame", "memory", "time", "stream", "dump", ""}

	for _, name := range names {
