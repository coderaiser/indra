//go:build ignore

package fixture

// wrong-return-type: Match does not return a Matcher, so it is not a plugin
// Match and must be left alone.
func Match() Replacer {
	return Replacer{`__a.Equal(__b, true)`: `__a.Ok(__b)`}
}
