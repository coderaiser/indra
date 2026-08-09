# remove-useless-comments

> Remove separator banner comments built from repeated `─` characters

## Example

### ❌ Incorrect

```go
package fixture

// ── imports ──────────────────────────────────────────────────────────────────

import "fmt"
```

### ✅ Correct

```go
package fixture

import "fmt"
```
