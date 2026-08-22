# indra

Rules for the indra linter itself.

## Rules

- ⛔ [remove-useless-match](#remove-useless-match)
- ⛔ [apply-exports-order](#apply-exports-order)

## Configuration

```toml
[match]
"*.go" = { "indra" = "on" }
```

## remove-useless-match

⛔ **disabled by default** — enable via `"indra" = "on"`.

Removes useless entries in a `Match()` map: a `nil` guard is reported as
useless, as is an empty `Matcher`.

### ❌ Incorrect

```go
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
package fixture

import . "coderaiser/indra/types"

func Report() string { return "remove Test.Skip call" }

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
```

## apply-exports-order

⛔ **disabled by default** — enable via `"indra" = "on"`.

Applies the canonical export order to plugin files. A file whose top-level
exported functions implement one of the known plugin shapes is reordered to
that shape's canonical order:

- Replacer: `Report`, `Match`, `Replace`
- Traverser: `Report`, `Fix`, `Traverse`
- Includer: `Report`, `Include`, `Fix`, `Filter`

Non-function declarations and unexported functions keep their positions;
exported functions outside the shape follow the shape functions.

### ❌ Incorrect

```go
package fixture

import . "coderaiser/indra/types"

func Report(_ Path) string { return "reorder" }

func Traverse() Traverser {
	return Traverser{}
}

func Fix(_ Path, _ map[string]any) {}
```

### ✅ Correct

```go
package fixture

import . "coderaiser/indra/types"

func Report(_ Path) string { return "reorder" }

func Fix(_ Path, _ map[string]any) {}

func Traverse() Traverser {
	return Traverser{}
}
```
