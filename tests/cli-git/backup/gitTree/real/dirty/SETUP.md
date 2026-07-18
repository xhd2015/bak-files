# Scenario

**Feature**: dirty gitTree backup copies worktree files and writes bak.stats

```
# operator
HOME=…/home  WORKING_ROLE=alice
git init ~/repo; commit README.md = "committed\n"; dirty → "dirty-readme\n"
config: "~/repo": { "gitTree": true }
operator -> bak-files backup --config bak.config.json
  -> files/HOME/repo/README.md == "dirty-readme\n"
  -> bak.stats[HOME/repo] hasChange=true, commitHash=HEAD
  -> exit 0
```

## Preconditions

- Local git repo with one commit then **uncommitted** edit to README.md.
- No remotes.
- Mapping `~` → `HOME`.

## Steps

1. `setupGitTreeBackup` with committed ≠ work content, dryRun=false.
2. Record BackupPath, StatsPath, MappingKey, CommitHash, Content.

## Context

- P5 exit criterion: dirty copies as designed + stats.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupGitTreeBackup(t, req, "committed\n", "dirty-readme\n", false)
	return nil
}
```
