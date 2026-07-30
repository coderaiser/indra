package indra_test

import (
	"testing"

	indra "coderaiser/indra"
	tape "github.com/coderaiser/go-tape"

	"github.com/lithammer/dedent"
)

func TestTrimIndent(t *testing.T) {
	tape.Test(t, "indra: TrimIndent removes common indent", func(t *tape.T) {
		result := indra.TrimIndent(dedent.Dedent(`
		if ok {
			return
		}
	`))
		t.Equal(result, dedent.Dedent(`
if ok {
	return
}
`))
		t.End()
	})

	tape.Test(t, "indra: TrimIndent preserves relative indent", func(t *tape.T) {
		result := indra.TrimIndent(dedent.Dedent(`
		if ok {
			if nested {
				return
			}
		}
	`))
		t.Equal(result, dedent.Dedent(`
if ok {
	if nested {
		return
	}
}
`))
		t.End()
	})

	tape.Test(t, "indra: TrimIndent ignores blank lines", func(t *tape.T) {
		result := indra.TrimIndent(dedent.Dedent(`

		if ok {

			return
		}

	`))
		t.Equal(result, dedent.Dedent(`

if ok {

	return
}

`))
		t.End()
	})

	tape.Test(t, "indra: TrimIndent returns empty string unchanged", func(t *tape.T) {
		t.Equal(indra.TrimIndent(""), "")
		t.End()
	})

	tape.Test(t, "indra: TrimIndent removes leading spaces", func(t *tape.T) {
		result := indra.TrimIndent("   hello\n\n   world")
		t.Equal(result, "hello\n\nworld")
		t.End()
	})
}
