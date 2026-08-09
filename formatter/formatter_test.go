package formatter_test

import (
	"testing"

	"coderaiser/indra/formatter"
	"fmt"
	"formjson"

	. "github.com/coderaiser/go-tape"
)

func TestChoose(t *testing.T) {
	Test(t, "formatter: exported: Choose returns a non-nil Func", func(t *T) {
		f := formatter.Choose()
		t.Ok(f)
		t.End()
	})

	Test(t, "formatter: exported: ChooseByName", func(t *T) {
		t.TB().Setenv("CI", "")
		f := formatter.ChooseByName("json")
		result := fmt.Sprintf("%p", f)
		t.Equal(result, fmt.Sprintf("%p", formjson.Format))
		t.End()
	})
}