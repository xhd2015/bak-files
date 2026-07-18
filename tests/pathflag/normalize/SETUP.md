# Scenario

**Feature**: input normalization before matching

```
# trim space, ToSlash, strip leading ./
caller -> Classify("  ./.bun  ") -> same as Classify(".bun")
```

## Preconditions

- Normalization does not change semantic match results for valid paths.

## Steps

1. Leaf sets a sugared path form.
2. Assert same attribute Result as the normalized catalog path.

## Context

- Invalid forms after normalize still error (see `invalid/`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
