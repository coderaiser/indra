//go:build ignore

package fixture

import foo "coderaiser/indra/internal/test"

var Test = indratest.For("some-rule", somePlugin{})
var y = someFunc()
var z *int