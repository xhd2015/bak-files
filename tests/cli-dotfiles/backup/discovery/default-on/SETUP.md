# Scenario

**Feature**: dots default-on — dry-run would handle `.bashrc`, skip `.cache`

```
# operator
HOME=…/home  files={}
write home/.bashrc = "export A=1\n"
write home/.cache/x = "cache-body\n"
operator -> bak-files backup --config … --dry-run
  -> exit 0
  -> targetDir unchanged (no writes)
  -> logs prefer would… for .bashrc and skip… for .cache
```

## Preconditions

- Empty files; includeDotFiles omitted (default true).
- `.cache` is pathflag Cache → DefaultSkipMask.

## Steps

1. setupDotsWorld emptyFilesConfig, backup, dryRun=true.
2. Plant `.bashrc` and `.cache/x`.
3. TargetDirBefore = tree fingerprint (absent).
4. Record BackupPath for bashrc mapping path (must not be created on dry-run).

## Context

- Requirement default-on scenario; dry-run for log + zero-write clarity.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", true)
	writeFile(t, filepath.Join(home, ".bashrc"), "export A=1\n")
	writeFile(t, filepath.Join(home, ".cache", "x"), "cache-body\n")
	req.Content = "export A=1\n"
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".cache", "x")
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	return nil
}
```
