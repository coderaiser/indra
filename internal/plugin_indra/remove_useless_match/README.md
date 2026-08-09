# remove-useless-match

> Remove useless entries in a `Match()` map

## Example

### ❌ Incorrect

```go
package fixture

import . "coderaiser/indra/types"

func Report() string { return "some plugin" }

func Match() Matcher {
    return Matcher{
        `Test.Skip(__a, __b, func(__a *Test.T) { __body })`: nil,
    }
}
```

### ✅ Correct

```go
package fixture

import . "coderaiser/indra/types"

func Report() string { return "some plugin" }

func Match() Matcher {
    return Matcher{}
}
```
