# Scenario

**Feature**: `.codex/sessions` are Codex session history (not cache)

```
caller -> Classify(".codex/sessions/2026/a.json")
  -> Rule=.codex/sessions, Flags=history, Owner=codex
```

## Preconditions

- Catalog: `.codex/sessions` → History, owner codex.
- Sessions must not be labeled cache or tmp.

## Steps

1. Set path under `.codex/sessions`.
2. Expect history flag and codex owner.

## Context

- Replaces the old owner-only expectation for session trees.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".codex/sessions/2026/a.json"
	return nil
}
```
