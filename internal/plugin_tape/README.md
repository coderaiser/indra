# tape

Rules for go-tape test files: 17 rules across five categories.

## Rules

### skips and onlys

- ✅ [remove-skip](#remove-skip)
- ✅ [remove-only](#remove-only)

### t.End()

- ✅ [add-t-end](#add-t-end)
- ✅ [remove-useless-t-end](#remove-useless-t-end)

### assertions

- ✅ [apply-assertions-order](#apply-assertions-order)
- ✅ [switch-expected-with-result](#switch-expected-with-result)
- ✅ [remove-default-messages](#remove-default-messages)
- ✅ [convert-deep-equal-to-equal](#convert-deep-equal-to-equal)

### operator conversions

- ✅ [convert-equal-to-deep-equal](#convert-equal-to-deep-equal)
- ✅ [convert-equal-to-ok](#convert-equal-to-ok)
- ✅ [convert-equal-to-not-ok](#convert-equal-to-not-ok)
- ✅ [convert-ok-to-not-ok](#convert-ok-to-not-ok)
- ✅ [convert-no-error-to-not-ok](#convert-no-error-to-not-ok)
- ✅ [extract-result-from-assertion](#extract-result-from-assertion)

### formatting

- ✅ [apply-dedent](#apply-dedent)
- ✅ [remove-useless-prefix](#remove-useless-prefix)
- ✅ [remove-useless-condition](#remove-useless-condition)

## Configuration

```toml
[match]
"*_test.go" = { "tape" = "on" }

[rules."tape/remove-skip"]
allowed = ["Suite"]
```

## remove-skip

Removes Test.Skip() calls from test functions. Extra allowed receivers come
from the `allowed` option.

### ❌ Incorrect

```go
Test.Skip(t, "foo: something", func(t *Test.T) {
	t.Equal(1, 1)
})
```

### ✅ Correct

```go
Test(t, "foo: something", func(t *Test.T) {
	t.Equal(1, 1)
})
```

## remove-only

Removes Test.Only() calls from test functions. Accepts the same `allowed`
option as remove-skip.

### ❌ Incorrect

```go
Test.Only(t, "foo: something", func(t *Test.T) {
	t.Equal(1, 1)
})
```

### ✅ Correct

```go
Test(t, "foo: something", func(t *Test.T) {
	t.Equal(1, 1)
})
```

## add-t-end

Adds a missing t.End() call at the end of a test function.

### ❌ Incorrect

```go
Test(t, "foo: something", func(t *T) {
	t.Equal(1, 1)
})
```

### ✅ Correct

```go
Test(t, "foo: something", func(t *T) {
	t.Equal(1, 1)
	t.End()
})
```

## remove-useless-t-end

Removes duplicate t.End() calls: every End beyond the first is a runtime
no-op.

### ❌ Incorrect

```go
func f(t *testing.T) {
	Test(t, "two ends", func(t *T) {
		t.Equal(1, 1)
		t.End()
		t.End()
	})
}
```

### ✅ Correct

```go
func f(t *testing.T) {
	Test(t, "two ends", func(t *T) {
		t.Equal(1, 1)
		t.End()
	})
}
```
