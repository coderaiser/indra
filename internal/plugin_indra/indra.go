package plugin_indra

import . "coderaiser/indra/types"

// Rules groups indra meta-rules.
// All rules are Off by default.
// Enable in .indra.toml: "indra" = "on"
var Rules = Nested{
	"remove-useless-match": Off("coderaiser/indra/internal/plugin_indra/remove_useless_match"),
}
