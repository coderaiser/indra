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
## apply-assertions-order

Orders assertions inside a tape test so they follow the canonical layout:
assertions come together, and `t.End()` is the last statement.

### ❌ Incorrect

```go
func TestFoo(t *testing.T) {
	t.Equal(1, 1)
	fmt.Println("gap")
	t.End()
}
```

### ✅ Correct

```go
func TestFoo(t *testing.T) {

	fmt.Println("gap")
	t.Equal(1, 1)

	t.End()
}
```

## switch-expected-with-result

Swaps the arguments of an assertion so `result` comes before `expected` — the
canonical `t.Equal(actual, expected)` shape.

### ❌ Incorrect

```go
func TestFoo(t *testing.T) {
	t.Equal(expected, actual)
	t.End()
}
```

### ✅ Correct

```go
func TestFoo(t *testing.T) {
	t.Equal(actual, expected)

	t.End()
}
```

## remove-default-messages

Removes assertion messages that repeat tape's built-in default ("should be
truthy", "should equal", ...) — they add noise doing nothing beyond the default.

### ❌ Incorrect

```go
func TestFoo(t *testing.T) {
	t.Ok(cond, "should be truthy")
	t.NotOk(other, "should be falsy")
	t.Equal(a, b, "should equal")
	t.DeepEqual(c, d, "should deep equal")
	t.End()
}
```

### ✅ Correct

```go
func TestFoo(t *testing.T) {
	t.Ok(cond)
	t.NotOk(other)
	t.Equal(a, b)
	t.DeepEqual(c, d)

	t.End()
}
```

## convert-deep-equal-to-equal

For primitive (non-slice, non-struct) arguments `DeepEqual` is the same as
`Equal` — the lighter `Equal` reads better.

### ❌ Incorrect

```go
func TestFoo(t *testing.T) {
	t.DeepEqual(name, "hello")
	t.End()
}
```

### ✅ Correct

```go
func TestFoo(t *testing.T) {
	t.Equal(name, "hello")

	t.End()
}
```
## convert-equal-to-deep-equal

For composite arguments (a slice literal here) `Equal` compares identity, not
content — `DeepEqual` compares the values the test actually means.

### ❌ Incorrect

```go
func f() {
	t.Equal(x, []Block{})
}
```

### ✅ Correct

```go
func f() {
	t.DeepEqual(x, []Block{})

}
```

## convert-equal-to-ok

`Equal(x, true)` is just `Ok(x)` — the shorter boolean assertion.

### ❌ Incorrect

```go
func f() {
	t.Equal(result, true)
}
```

### ✅ Correct

```go
func f() {
	t.Ok(result)

}
```

## convert-equal-to-not-ok

`Equal(x, nil)` (or `Equal(x, false)`) is `NotOk(x)`.

### ❌ Incorrect

```go
func f() {
	t.Equal(err, nil)
}
```

### ✅ Correct

```go
func f() {
	t.NotOk(err)

}
```

## convert-ok-to-not-ok

A negated `Ok` — `Ok(x == nil)` or `Ok(!x)` — reads more directly as `NotOk`.

### ❌ Incorrect

```go
func f() {
	t.Ok(err == nil)
}
```

### ✅ Correct

```go
func f() {
	t.NotOk(err)

}
```

## convert-no-error-to-not-ok

`NoError(err)` is the same assertion as `NotOk(err)`.

### ❌ Incorrect

```go
func f() {
	t.NoError(err)
}
```

### ✅ Correct

```go
func f() {
	t.NotOk(err)

}
```

## extract-result-from-assertion

Pulls an inline expression out of an assertion into its own variable so the
result is named and the assertion is easier to read.

### ❌ Incorrect

```go
func f() {
	t.DeepEqual(someFunc(a, b), expected)
}
```

### ✅ Correct

```go
func f() {
	result := someFunc(a, b)
	t.DeepEqual(result, expected)

}
```
## apply-dedent

Removes the `dedent.Dedent(...)` wrapper around a raw string literal, leaving
the literal in place for the formatter to normalize.

### ❌ Incorrect

```go
func f() []byte {
	return []byte(dedent.Dedent(`
		[ignore]
		patterns = ["vendor/**"]
		`))
}
```

### ✅ Correct

```go
func f() []byte {
	return []byte(`
		[ignore]
		patterns = ["vendor/**"]
		`)
}
```

## remove-useless-prefix

Drops a redundant receiver prefix — `tape.Test(...)` used inside a file that
already dot-imports the package reads as plain `Test(...)`.

### ❌ Incorrect

```go
func TestFoo(t *testing.T) {
	tape.Test(t, "foo: bar", func(t *tape.T) {
		t.Equal(1, 1)
		t.End()
	})
}
```

### ✅ Correct

```go
func TestFoo(t *testing.T) {
	Test(t, "foo: bar", func(t *T) {
		t.Equal(1, 1)
		t.End()
	})
}
```

## remove-useless-condition

Simplifies assertions whose condition is redundant: `Ok(x != nil)` means the
same as `Ok(x)`, and the easier form reads better. Message arguments are kept.

### ❌ Incorrect

```go
func f() {
	t.Ok(err != nil, "err should be set")
}
```

### ✅ Correct

```go
func f() {
	t.Ok(err, "err should be set")

}
```
