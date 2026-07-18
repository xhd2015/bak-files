# Scenario

**Feature**: ordinary home-relative file has no attributes

```
caller -> Classify("Documents/notes.txt") -> zero Result, nil error
```

## Preconditions

- Path does not match catalog, segments, log suffix, or owner prefixes.

## Steps

1. Set `RelPath` to `Documents/notes.txt`.
2. Expect empty Result fields.

## Context

- Baseline negative case for the matcher.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "Documents/notes.txt"
	return nil
}
```
