# Scenario

**Feature**: `.local/share/claude` is Claude install/version cache (multi-flag)

```
caller -> Classify(".local/share/claude/versions/2.1.198")
  -> Rule=.local/share/claude, Flags=cache|binary, Owner empty (or claude if added)
```

## Preconditions

- Catalog: `.local/share/claude` → Cache|Binary.
- Owner may be empty unless `OwnerClaude` is introduced; this leaf allows empty.
- Does not apply to `.local/bin/claude`.

## Steps

1. Set path under `.local/share/claude/versions`.
2. Expect `cache|binary` in ascending bit order.

## Context

- Bit order: cache before binary (not `binary|cache`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/share/claude/versions/2.1.198"
	return nil
}
```
