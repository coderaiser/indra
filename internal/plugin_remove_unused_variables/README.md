# remove-unused-variables

> Remove unused imports, private functions, consts, and variables

## Rules

- ✅ [remove-unused-import](#remove-unused-import)
- ✅ [remove-unused-variable](#remove-unused-variable)

## remove-unused-import

Removes an unused import spec.

### ❌ Incorrect

```go
package fixture

import (
    "fmt"
    "strings"
)

func f() {
    fmt.Println("hi")
}
```

### ✅ Correct

```go
package fixture

import (
    "fmt"
)

func f() {
    fmt.Println("hi")
}
```

## remove-unused-variable

Removes an unused const, variable, or private function declaration.

### ❌ Incorrect

```go
package fixture

const used = 1
const unused = 2

var x = 1
var y = 2

func usedFunc() {}
func helper() {}

func f() {
    println(used, x, usedFunc)
}

func g() {}
```

### ✅ Correct

```go
package fixture

const used = 1

var x = 1

func usedFunc() {}

func f() {
    println(used, x, usedFunc)
}
```
