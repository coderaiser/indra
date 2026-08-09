# convert-no-error-to-not-ok

> Convert `NoError(err)` to `NotOk(err)`

## Example

### ❌ Incorrect

```go
package fixture

import tape "github.com/coderaiser/go-tape"

func f(t *tape.T) {
    x := 1
    foo()
    t.NotOk(other)
    t.NoError(err)
}
```

### ✅ Correct

```go
package fixture

import tape "github.com/coderaiser/go-tape"

func f(t *tape.T) {
    x := 1
    foo()
    t.NotOk(other)
    t.NotOk(err)
}
```
