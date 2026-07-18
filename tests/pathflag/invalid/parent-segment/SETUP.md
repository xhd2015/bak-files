# Scenario

**Feature**: any `..` path segment is rejected

```
# parent directory traversal in relative path
caller -> Classify("foo/../.cache") -> error
```

## Preconditions

- After normalize, any segment equal to `..` is illegal.

## Steps

1. Set `RelPath` containing a `..` segment.
2. Classify must error.

## Context

- Protects callers from accidental path escape; no silent cleaning of `..`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "foo/../.cache"
	return nil
}
```
