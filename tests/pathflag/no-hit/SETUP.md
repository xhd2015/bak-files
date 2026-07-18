# Scenario

**Feature**: unclassified relative paths return zero Result

```
# no catalog rule, no owner prefix
caller -> Classify("Documents/notes.txt") -> Result{} (nil error)
# negatives: paths near catalog rules that must not inherit skip flags
```

## Preconditions

- Path is valid (relative, no `..`) but unmatched, or is a deliberate negative
  against over-broad catalog prefixes.

## Steps

1. Leaf sets a path that must not receive skip attributes.
2. Assert zero Result and empty Err.

## Context

- Zero Result means empty Rule, Reason, Flags string, and Owner.
- Negatives: V2rayU config (not v2ray-core), `.local/state/...`, `.local/bin/claude`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
