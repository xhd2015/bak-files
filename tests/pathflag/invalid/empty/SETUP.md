# Scenario

**Feature**: empty path after normalize is an error

```
# empty or whitespace-only string
caller -> Classify("") | Classify("   ") -> error
```

## Preconditions

- Empty after trim is rejected.

## Steps

1. Set `RelPath` to empty string.
2. Classify must error.

## Context

- Whitespace-only is equivalent after trim (covered by same empty rule).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ""
	return nil
}
```
