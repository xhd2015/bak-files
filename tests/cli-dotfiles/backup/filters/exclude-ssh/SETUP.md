# Scenario

**Feature**: `--exclude .ssh` skips secrets that would otherwise auto-backup

```
# operator
write home/.bashrc = "rc\n"
write home/.ssh/id_rsa = "secret\n"
operator -> bak-files backup --config … --exclude .ssh
  -> .bashrc backed up; .ssh/id_rsa MUST NOT be under targetDir
  -> exit 0
```

## Preconditions

- Dots on; `.ssh` included by default (not pathflag-skipped).

## Steps

1. setupDotsWorld with `--exclude .ssh`.
2. Plant `.bashrc` and `.ssh/id_rsa`.
3. KeepBackupPath = bashrc; ExcludedBackupPath = id_rsa under mapping.

## Context

- Locked decision: secrets included by default until exclude/dotExcludes.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false,
		"--exclude", ".ssh")
	writeFile(t, filepath.Join(home, ".bashrc"), "rc\n")
	writeFile(t, filepath.Join(home, ".ssh", "id_rsa"), "secret\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.KeepContent = "rc\n"
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".ssh", "id_rsa")
	return nil
}
```
