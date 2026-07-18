# Scenario

**Feature**: match priority — longest path rule, segments, log suffix

```
# path catalog longest-prefix first
# else segment rules (node_modules, upload-chunks)
# else basename .log suffix
caller -> Classify(path) -> winning Rule
```

## Preconditions

- Leaves isolate one priority or fallback rule at a time.

## Steps

1. Leaf sets a path that exercises longest-prefix, segment, or suffix.
2. Assert the winning Rule and flags.

## Context

- Generic binary content rule is out of scope for Classify.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
