# Scenario

**Feature**: restore dry-run respects pathflag skip for `.codex/.tmp`

```
# operator
seed files/HOME/alice/.codex/config = "cfg\n"
seed files/HOME/alice/.codex/.tmp/plugin = "tmp\n"
dest home absent for both
operator -> bak-files restore --config … --dry-run
  -> exit 0
  -> dest paths still absent
  -> prefer skip log for .tmp / pathflag reason
```

## Preconditions

- Empty files config still used so discovery maps `~/.codex` symmetrically;
  backup store pre-seeded under mapping layout.
- Destinations not pre-created.

## Steps

1. setupDotsWorld emptyFilesConfig, restore, dryRun=true.
2. Seed BackupPath keep + excluded under TargetDir.
3. DestPath = home/.codex/.tmp/plugin; DestBefore empty.
4. TargetDirBefore not required for restore dest checks.

## Context

- Dry-run must not write dest; skip log proves policy applied on restore path.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "restore", true)
	// Seed backup store as if a prior backup included both (policy may still skip restore of tmp).
	writeFile(t, mappingBackup(req.TargetDir, "alice", ".codex", "config"), "cfg\n")
	writeFile(t, mappingBackup(req.TargetDir, "alice", ".codex", ".tmp", "plugin"), "tmp\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".codex", "config")
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".codex", ".tmp", "plugin")
	req.DestPath = filepath.Join(home, ".codex", ".tmp", "plugin")
	req.SourcePath = filepath.Join(home, ".codex", "config")
	req.DestBefore = readFileOrEmpty(req.DestPath)
	req.Content = "cfg\n"
	return nil
}
```
