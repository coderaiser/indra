# tape

Rules for go-tape test files.

## Rules

- ✅ [remove-skip](#remove-skip)
- ✅ [add-t-end](#add-t-end)
- ✅ [convert-equal-to-deep-equal](#convert-equal-to-deep-equal)
- ✅ [convert-equal-to-not-ok](#convert-equal-to-not-ok)
- ✅ [convert-ok-to-not-ok](#convert-ok-to-not-ok)
- ✅ [extract-result-from-assertion](#extract-result-from-assertion)
- ✅ [remove-useless-condition](#remove-useless-condition)
- ✅ [remove-useless-prefix](#remove-useless-prefix)
- ✅ [convert-no-error-to-not-ok](#convert-no-error-to-not-ok)

## Configuration

```toml
[match]
"*_test.go" = { "tape" = "on" }
```

## remove-skip

Removes t.Skip() calls from test functions.

### ❌ Incorrect

```go
package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test.Skip(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})
}
```

### ✅ Correct

```go
package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})

}
```

## add-t-end

Adds a missing t.End() call at the end of a test function.

### ❌ Incorrect

```go
package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
	})
}
```

### ✅ Correct

```go
package fixture

import (
	Test "github.com/coderaiser/go-tape"
	"testing"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: something", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})

}
```

## convert-equal-to-deep-equal

Uses DeepEqual instead of Equal when comparing slices.

### ❌ Incorrect

```go
package fixture

// convert-equal-to-deep-equal is the canonical happy path: Equal used on a slice.
func f() {
	t.Equal(x, []Block{})
}
```

### ✅ Correct

```go
package fixture

// convert-equal-to-deep-equal is the canonical happy path: Equal used on a slice.
func f() {
	t.DeepEqual(x, []Block{})

}
```

## convert-equal-to-not-ok

Converts Equal(err, nil) to NotOk(err).

### ❌ Incorrect

```go
package fixture

func f() {
	t.Equal(err, nil)
}
```

### ✅ Correct

```go
package fixture

func f() {
	t.NotOk(err)

}
```

## convert-ok-to-not-ok

Converts Ok(err == nil) and Ok(!err) to NotOk(err).

### ❌ Incorrect

```go
package fixture

func f() {
	t.Ok(err == nil)
	t.Ok(!err)
}
```

### ✅ Correct

```go
package fixture

func f() {
	t.NotOk(err)
	t.NotOk(err)
}
```

## extract-result-from-assertion

Extracts an inline expression from an assertion into a named result variable.

### ❌ Incorrect

```go
package fixture

func f() {
	t.DeepEqual(someFunc(a, b), expected)
}
```

### ✅ Correct

```go
package fixture

func f() {
	result := someFunc(a, b)
	t.DeepEqual(result, expected)

}
```

## remove-useless-condition

Removes a useless condition inside Ok().

### ❌ Incorrect

```go
package fixture

func f() {
	t.Ok(err != nil)
}
```

### ✅ Correct

```go
package fixture

func f() {
	t.Ok(err)

}
```

## remove-useless-prefix

Removes the redundant tape prefix so the tape package is dot-imported.

### ❌ Incorrect

```go
package fixture

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	tape.Test(t, "foo: bar", func(t *tape.T) {
		t.Equal(1, 1)
		t.End()
	})
}
```

### ✅ Correct

```go
package fixture

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestFoo(t *testing.T) {
	Test(t, "foo: bar", func(t *T) {
		t.Equal(1, 1)
		t.End()
	})
}
```

## convert-no-error-to-not-ok

Converts NoError(err) to NotOk(err).

### ❌ Incorrect

```go
package fixture

import tape "github.com/coderaiser/go-tape"

func f() {
	x := 1
	foo()
	t.NotOk(other)
	t.NoError(err)
}
```

### ✅ Correct

```go
package fixture

import tape "github.com/coderaiser/go-tape"

func f() {
	x := 1
	foo()
	t.NotOk(other)
	t.NotOk(err)
}
```
