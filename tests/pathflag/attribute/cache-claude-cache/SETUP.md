# Scenario

**Feature**: `.claude/cache` is Claude tool cache

```
caller -> Classify(".claude/cache/x")
  -> Rule=.claude/cache, Flags=cache, Owner=claude
```

## Preconditions

- Catalog: `.claude/cache` → Cache, owner claude.
- Whole `.claude` root is not a catalog rule.

## Steps

1. Set path under `.claude/cache`.
2. Expect cache flags and owner claude.

## Context

- Fine prefix only; regenerable Claude cache data.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".claude/cache/x"
	return nil
}
```
