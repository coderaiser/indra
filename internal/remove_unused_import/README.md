# remove-unused-import

> Remove unused imports

## Rule

```toml
[rules]
"remove-unused-import" = "on"
```

### ❌ Incorrect

```go
//go:build ignore

package fixture

import "fmt"

func f() {}
```

### ✅ Correct

```go
//go:build ignore

package fixture

func f() {}
```
