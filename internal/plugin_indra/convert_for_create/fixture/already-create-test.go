//go:build ignore

package fixture

import . "coderaiser/indra/internal/test"

var Test = createTest("some-rule", somePlugin{})
