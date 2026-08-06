//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("some-rule", somePlugin{})
