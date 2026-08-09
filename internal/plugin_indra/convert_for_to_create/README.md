# convert-for-to-create-test

> Convert `indratest.For(...)` calls to `CreateTest(...)`

## Example

### ❌ Incorrect

```go
package fixture

import (
    indratest "coderaiser/indra/internal/test"
)

func TestFoo(t *testing.T) {
    indratest.For(t, "foo: bar", func(t *indratest.T) {
        t.Equal(1, 1)
    })
}
```

### ✅ Correct

```go
package fixture

import (
    . "coderaiser/indra/internal/test"
)

func TestFoo(t *testing.T) {
    CreateTest(t, "foo: bar", func(t *T) {
        t.Equal(1, 1)
    })
}
```
