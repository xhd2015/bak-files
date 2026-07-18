# Scenario

**Feature**: Classify rejects invalid normalized paths

```
# empty, absolute, or ".." segments
caller -> pathflag.Classify(badPath) -> error (no Result attributes)
```

## Preconditions

- Paths are rejected after normalize (trim, ToSlash, strip `./`).

## Steps

1. Leaf sets an invalid `RelPath`.
2. Run Classify; Assert expects non-empty `Response.Err`.

## Context

- Invalid inputs must not return a silent zero Result.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
