# replace-test-message

> Include fixture name in test message

## Example

### ❌ Incorrect

```go
import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *test.T) {
    Test(t, "remove-skip: report", func(t *T) {
        t.Report("remove-skip", "remove Test.Skip call")
        t.End()
    })
}
```

### ✅ Correct

```go
import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *test.T) {
    Test(t, "remove-skip: report: remove-skip", func(t *T) {
        t.Report("remove-skip", "remove Test.Skip call")
        t.End()
    })
}
```

## Configuration

```toml
[match]
"*.go" = { "indra/replace-test-message" = "off" }
```