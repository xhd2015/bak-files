# Scenario

**Feature**: dirty gitPatch + `--dry-run` → no patch, exit 0

```
# operator
dirty ~/repo; config gitPatch:true
operator -> bak-files backup --config … --dry-run
  -> no *.patch under files/
  -> prefer dry-run/would logs
  -> exit 0
```

## Preconditions

- Dirty worktree; dry-run true.
- TargetDirBefore fingerprint empty.

## Steps

1. `setupGitPatchBackup` dryRun=true.

## Context

- Ensures dry-run short-circuits patch generation.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupGitPatchBackup(t, req, "base-line\n", "patched-line\n", true)
	return nil
}
```
