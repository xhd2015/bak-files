# Scenario

**Feature**: dirty gitPatch backup produces a patch file under mapping path

```
# operator
commit README "base-line\n"; dirty → "patched-line\n"
config: "~/repo": { "gitPatch": true }
operator -> bak-files backup --config bak.config.json
  -> under files/HOME/repo/ at least one *.patch (or BackupPath)
  -> patch body looks like git/unified diff
  -> exit 0
```

## Preconditions

- Dirty worktree vs HEAD.
- No remotes; base = HEAD for MVP.

## Steps

1. `setupGitPatchBackup` with distinct committed vs work content, dryRun=false.

## Context

- P5 exit criterion: gitPatch produces patch when dirty.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupGitPatchBackup(t, req, "base-line\n", "patched-line\n", false)
	return nil
}
```
