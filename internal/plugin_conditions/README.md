# conditions

Rules that apply to any assertion library. Nothing here is tape-specific, so
the group is enabled for all files by default — no glob restriction needed.

## Rules

- ✅ [remove-useless-comments](#remove-useless-comments)

## remove-useless-comments

Removes separator banner comments built from repeated `─` characters.

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