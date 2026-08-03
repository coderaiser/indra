package plugin_indra_test

import (
	"testing"

	. "coderaiser/indra/internal/plugin-indra"
	. "coderaiser/indra/types"
	tape "github.com/coderaiser/go-tape"
)

func TestRules(t *testing.T) {
	tape.Test(t, "plugin-indra: remove-useless-match is present", func(t *tape.T) {
		_, ok := Rules["remove-useless-match"]
		t.Equal(ok, true)
		t.End()
	})

	tape.Test(t, "plugin-indra: remove-useless-match is Off by default", func(t *tape.T) {
		entry := Rules["remove-useless-match"].(PluginEntry)
		t.Equal(entry.Enabled, false)
		t.End()
	})
}
