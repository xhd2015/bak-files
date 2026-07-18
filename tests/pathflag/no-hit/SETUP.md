# Scenario

**Feature**: unclassified relative paths return zero Result

```
# no catalog rule, no owner prefix
caller -> Classify("Documents/notes.txt") -> Result{} (nil error)
```

## Preconditions

- Path is valid (relative, no `..`) but unmatched.

## Steps

1. Leaf sets an ordinary relative path.
2. Assert zero Result and empty Err.

## Context

- Zero Result means empty Rule, Reason, Flags string, and Owner.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
