//go:build ignore

package fixture

import indratest "coderaiser/indra/internal/test"

var Test = indratest.For("some-rule", somePlugin{})

func f(t *indratest.T) {}