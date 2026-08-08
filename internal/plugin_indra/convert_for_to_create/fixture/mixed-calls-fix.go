//go:build ignore

package fixture

import foo "coderaiser/indra/internal/test"

var Test = CreateTest("some-rule", somePlugin{})
var y = someFunc()
var z *int
