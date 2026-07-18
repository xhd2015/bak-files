# Scenario

**Feature**: real gitTree backup (no `--dry-run`) may write targetDir and bak.stats

```
# real run
operator -> bak-files backup --config …   # no --dry-run
  -> may write files/ and bak.stats
```

## Preconditions

- DryRun is false.
- May mutate WorkDir/files and WorkDir/bak.stats.

## Steps

1. Leaves choose dirty vs clean worktree.

## Context

- Real mode is the only path that should persist bak.stats for gitTree.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
