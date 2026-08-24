package plugin_package_json_test

import (
	"testing"

	. "coderaiser/indra/internal/plugin_package_json"

	. "github.com/coderaiser/go-tape"
)

func TestRules(t *testing.T) {
	Test(t, "package-json: Rules returns 1 rule", func(t *T) {
		result := len(Rules())
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "package-json: rename-version is registered", func(t *T) {
		t.Equal(Rules()[0].Name, "rename-version")
		t.End()
	})
}
