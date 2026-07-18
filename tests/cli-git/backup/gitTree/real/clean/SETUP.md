# Scenario

**Feature**: clean gitTree backup skips copy (INFO); no targetDir churn

```
# operator
git init ~/repo; commit README.md = "clean-committed\n"; worktree clean
config: "~/repo": { "gitTree": true }
operator -> bak-files backup --config bak.config.json
  -> no files/HOME/repo/README.md (or targetDir fingerprint unchanged)
  -> bak.stats absent or no hasChange:true for mapping after skip
  -> prefer INFO / skip / clean in logs
  -> exit 0
```

## Preconditions

- Local git repo with committed README only; **no** uncommitted changes.
- No remotes.

## Steps

1. `setupGitTreeBackup` with committed == work content, dryRun=false.
2. Capture TargetDirBefore / StatsBefore (empty).

## Context

- P5 exit criterion: clean tree skips backup.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupGitTreeBackup(t, req, "clean-committed\n", "clean-committed\n", false)
	return nil
}
```
