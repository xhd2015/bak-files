# Scenario

**Feature**: `.claude/backups` is Claude local backup cache

```
caller -> Classify(".claude/backups/.claude.json.backup.1")
  -> Rule=.claude/backups, Flags=cache, Owner=claude
```

## Preconditions

- Catalog: `.claude/backups` → Cache, owner claude.
- Whole `.claude` root is not a catalog rule.

## Steps

1. Set path under `.claude/backups`.
2. Expect cache flags and owner claude.

## Context

- Only listed Claude prefixes are cataloged, not the entire `.claude` tree.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".claude/backups/.claude.json.backup.1"
	return nil
}
```
