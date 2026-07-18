# Scenario

**Feature**: real gitPatch backup may write a patch under targetDir

```
# real
operator -> bak-files backup --config …   # no --dry-run
  -> may write *.patch under mapping
```

## Preconditions

- DryRun is false.

## Steps

1. Leaf sets dirty worktree.

## Context

- Real mode only path that should create patch artifacts.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
