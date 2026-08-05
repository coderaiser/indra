# remove-unused-import

> Remove unused imports

## Example

### ❌ Incorrect

```go
import "fmt"

func f() {}
```

### ✅ Correct

```go
func f() {}
```

## Disable

```toml
[rules]
"remove-unused-import" = "off"
```
