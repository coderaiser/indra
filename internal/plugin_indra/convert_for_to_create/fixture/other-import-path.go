//go:build ignore

package fixture

import other "other/path"

var Test = indratest.For("some-rule", somePlugin{})