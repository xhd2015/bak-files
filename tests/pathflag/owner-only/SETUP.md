# Scenario

**Feature**: owner prefix without attribute rule

```
# owner table longest-prefix; Flags stay zero, no Rule
caller -> Classify(".codex/sessions/x") -> Owner=codex, Flags="", Rule=""
```

## Preconditions

- Path matches an owner prefix but not an attribute catalog rule.

## Steps

1. Leaf sets a path under an owner root that has no catalog entry.
2. Assert Owner only.

## Context

- Attribute rules that set Owner must still agree with the owner table; these
  leaves are the pure owner-only branch.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
