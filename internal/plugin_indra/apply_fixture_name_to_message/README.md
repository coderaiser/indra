# apply-fixture-name-to-message

> Prefix test message with rule name

## Example

### ❌ Incorrect

```go
import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *test.T) {
    Test(t, "report Test.Skip call", func(t *T) {
        t.End()
    })
}
```

### ✅ Correct

```go
import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *test.T) {
    Test(t, "remove-skip: report Test.Skip call", func(t *T) {
        t.End()
    })
}
```

## Configuration

```toml
[match]
"*.go" = { "indra/apply-fixture-name-to-message" = "off" }
```