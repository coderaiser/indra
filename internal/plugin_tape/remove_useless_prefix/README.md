# remove-useless-prefix

> Remove the redundant tape prefix so the tape package is dot-imported

## Example

### ❌ Incorrect

```go
package fixture

import (
    "testing"

    tape "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
    tape.Test(t, "foo: bar", func(t *tape.T) {
        t.Equal(1, 1)
        t.End()
    })
}
```

### ✅ Correct

```go
package fixture

import (
    "testing"

    . "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
    Test(t, "foo: bar", func(t *T) {
        t.Equal(1, 1)
        t.End()
    })
}
```
