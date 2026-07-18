# Scenario

**Feature**: `.claude/projects` is Claude per-workspace project cache

```
caller -> Classify(".claude/projects/-Users-x/y")
  -> Rule=.claude/projects, Flags=cache, Owner=claude
```

## Preconditions

- Catalog: `.claude/projects` → Cache, owner claude.
- Whole `.claude` root is not a catalog rule.

## Steps

1. Set path under `.claude/projects`.
2. Expect cache flags and owner claude.

## Context

- Workspace-local Claude project state; regenerable, not history.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".claude/projects/-Users-x/y"
	return nil
}
```
