# Scenario

**Feature**: gitPatch with `--dry-run` never writes a patch file

```
# dry-run
operator -> bak-files backup --config … --dry-run
  -> no *.patch under targetDir
```

## Preconditions

- Args include `--dry-run`.

## Steps

1. Leaf uses dirty worktree so real mode would write a patch.

## Context

- P5 exit criterion: dry-run gitPatch no patch file.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
