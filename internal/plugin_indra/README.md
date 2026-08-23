# indra

Rules for the indra linter itself.

## Rules

- ⛔ [remove-useless-match](#remove-useless-match)
- ⛔ [apply-exports-order](#apply-exports-order)
- ⛔ [convert-for-to-create-test](#convert-for-to-create-test)
- ⛔ [replace-test-message](#replace-test-message)
- ⛔ [apply-fixture-name-to-message](#apply-fixture-name-to-message)
- ⛔ [convert-inspect-to-traverse](#convert-inspect-to-traverse)
- ⛔ [apply-compare](#apply-compare)

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

## convert-for-to-create-test

⛔ **disabled by default** — enable via `"indra" = "on"`.

Converts the older `indratest.For(...)` test helper into the exported
`CreateTest`, importing it with a dot-import so the `go-tape` test DSL reads the
same way as every other rule.

### ❌ Incorrect

```go
package fixture

import indratest "coderaiser/indra/internal/test"

var Test = indratest.For("some-rule", somePlugin{})
```

### ✅ Correct

```go
package fixture

import . "coderaiser/indra/internal/test"

var Test = CreateTest("some-rule", somePlugin{})
```

## replace-test-message

⛔ **disabled by default** — enable via `"indra" = "on"`.

Replaces a literal message in a `Test(t, ...)` call with the rule's canonical
name, so the reported message always matches the rule being tested.

### ❌ Incorrect

```go
func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})
}
```

### ✅ Correct

```go
func f(t *testing.T) {
	Test(t, "remove-skip: report: remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})
}
```

## apply-fixture-name-to-message

⛔ **disabled by default** — enable via `"indra" = "on"`.

Syncs the fixture name in a `Test(t, "...: report: <name>")` message so it
points at the fixture actually exercised by the test.

### ❌ Incorrect

```go
func f(t *testing.T) {
	Test(t, "remove-skip: report: wrong-name", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})
}
```

### ✅ Correct

```go
func f(t *testing.T) {
	Test(t, "remove-skip: report: remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})
}
```

## convert-inspect-to-traverse

⛔ **disabled by default** — enable via `"indra" = "on"`.

Flags `ast.Inspect` calls in an AST-walking plugin file (one declaring a
`Traverse() Traverser`) so the file can move to indra's `Traverse` API. This
rule is report-only; it does not rewrite the call.

### ❌ Incorrect

```go
func find(p Path, push func(Path)) {
	ast.Inspect(p.Node, func(n ast.Node) bool { return true })
}
```

### ✅ Correct

```go
func find(p Path, push func(Path)) {
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			// visit nodes via the path tree
		},
	})
}
```

## apply-compare

⛔ **disabled by default** — enabled via `"indra" = "on"`.

Replaces the verbose `compare.GetTemplateValues` idiom in a plugin guard with
the compact `operator.Compare` helper.

### ❌ Incorrect

```go
import (
	"go/ast"

	"coderaiser/indra/compare"
)

func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if compare.GetTemplateValues(s, "__.End()") != nil {
			return true
		}
	}
	return false
}
```

### ✅ Correct

```go
import (
	. "coderaiser/indra/operator"
	"go/ast"
)

func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if Compare(s, "__.End()") {
			return true
		}
	}
	return false
}
```
