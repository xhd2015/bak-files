# Scenario

**Feature**: absolute paths are rejected

```
# leading slash (unix absolute)
caller -> Classify("/home/user/.cache") -> error
```

## Preconditions

- Absolute paths are not home-relative.

## Steps

1. Set `RelPath` to a path with leading `/`.
2. Classify must error.

## Context

- Windows drive-absolute forms are also rejected by the same absolute rule
  (this leaf covers the unix `/` case).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "/home/user/.cache"
	return nil
}
```
