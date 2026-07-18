# Scenario

**Feature**: gitTree backup with `--dry-run` never writes targetDir or bak.stats

```
# dry-run
operator -> bak-files backup --config … --dry-run
  -> log would…; no targetDir; bak.stats unchanged
```

## Preconditions

- Args include `--dry-run`.
- Zero mutation of backup artifacts.

## Steps

1. Leaf uses a **dirty** worktree so real mode would have written; dry-run must still no-op writes.

## Context

- Proves dry-run short-circuits gitTree write path (and must not touch remotes).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
