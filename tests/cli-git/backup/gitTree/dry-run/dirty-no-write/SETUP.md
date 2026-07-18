# Scenario

**Feature**: dirty gitTree + `--dry-run` → zero writes, exit 0

```
# operator
dirty ~/repo README; config gitTree:true
operator -> bak-files backup --config … --dry-run
  -> no files/HOME/repo/README.md
  -> bak.stats unchanged (absent)
  -> prefer dry-run/would in logs
  -> exit 0
```

## Preconditions

- Dirty worktree (would backup if real).
- Fingerprints TargetDirBefore / StatsBefore captured before Run.

## Steps

1. `setupGitTreeBackup` committed≠work, dryRun=true.

## Context

- P5 exit criterion: dry-run does not write patch/stats (here: tree + stats).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupGitTreeBackup(t, req, "committed\n", "would-be-backed\n", true)
	return nil
}
```
