# indra

> Rules for the indra linter itself.

## Rules

- ⛔ [remove-useless-match](#remove-useless-match)

## Configuration

```toml
[match]
"*.go" = { "indra" = "on" }
```

## remove-useless-match

⛔ **disabled by default** — enable via `"indra" = "on"`.

> Removes useless entries in a `Match()` map: a `nil` guard is reported as
> useless, as is an empty `Matcher`.

### ❌ Incorrect

```go
//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report() string { return "remove Test.Skip call" }

func Match() Matcher {
	return Matcher{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: nil,
	}
}

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
```

### ✅ Correct

```go
//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report() string { return "remove Test.Skip call" }

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
```
