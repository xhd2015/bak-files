# Scenario

**Feature**: `.claude/downloads` is Claude downloads cache

```
caller -> Classify(".claude/downloads/x")
  -> Rule=.claude/downloads, Flags=cache, Owner=claude
```

## Preconditions

- Catalog: `.claude/downloads` → Cache, owner claude.
- Whole `.claude` root is not a catalog rule.

## Steps

1. Set path under `.claude/downloads`.
2. Expect cache flags and owner claude.

## Context

- Fine prefix only; download staging is reclaimable.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".claude/downloads/x"
	return nil
}
```
