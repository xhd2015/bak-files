# Scenario

**Feature**: attribute catalog path rules set Rule, Reason, Flags, Owner

```
# longest matching path rule among catalog entries
caller -> Classify(catalogPath)
  -> Rule, Reason, Flags (bit names), Owner
```

## Preconditions

- Path equals a catalog rule or is under `rule + "/"`.

## Steps

1. Leaf sets a representative catalog path (exact or nested under rule).
2. Assert Rule, Reason, Flags string, Owner.

## Context

- Representative rows cover Cache, Logs, Tmp, Binary, Trash, Meta, History,
  and multi-flag combos. Session trees use `history` (never cache/tmp).
- `.local` is fine-prefix only; whole `.local` is not a rule.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	return nil
}
```
