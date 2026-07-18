# Scenario

**Feature**: leading/trailing spaces are trimmed

```
caller -> Classify("  .cache  ")
  -> Rule=.cache, Flags=cache, Owner empty
```

## Preconditions

- Trim space is the first normalize step.

## Steps

1. Set path with surrounding spaces.
2. Expect `.cache` attribute hit.

## Context

- After trim, empty string still errors (see `invalid/empty`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "  .cache  "
	return nil
}
```
