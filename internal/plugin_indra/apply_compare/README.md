# apply-compare

> Use `Compare` instead of `compare.GetTemplateValues(...) != nil`

## Example

### ❌ Incorrect

```go
package apply_compare

import (
    "coderaiser/indra/compare"
    "coderaiser/indra/types"
)

func Match() Matcher {
    return Matcher{
        `Test.Skip(__a, __b, func(__a *Test.T) { __body })`: func(v Vars, b *ast.BlockStmt) bool {
            if compare.GetTemplateValues(v.Args, b) != nil {
                return true
            }
            return false
        },
    }
}
```

### ✅ Correct

```go
package apply_compare

import (
    "coderaiser/indra/compare"
    "coderaiser/indra/types"
)

func Match() Matcher {
    return Matcher{
        `Test.Skip(__a, __b, func(__a *Test.T) { __body })`: func(v Vars, b *ast.BlockStmt) bool {
            if Compare(v.Args, b) {
                return true
            }
            return false
        },
    }
}
```
