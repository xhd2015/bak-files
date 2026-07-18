# Scenario

**Feature**: after access-denied warn, other good files still backed up

```
# operator
HOME=…/home  files={}
write home/.bashrc = "export A=1\n"
write home/.other = "other-body\n"
mkdir home/.private ; write secret ; chmod 000
operator -> bak-files backup --config …
  -> exit 0
  -> warning skip for .private
  -> both .bashrc and .other present under mapping with correct bodies
```

## Preconditions

- Same unreadable-dir fixture as unreadable-dir-warn, plus second good file
  `.other` (top-level dot, not pathflag-skipped).
- Empty files; default includeDotFiles.
- Cleanup: chmod restore via makeUnreadableDir Cleanup + Assert defer.

## Steps

1. setupDotsWorld emptyFilesConfig, real backup.
2. Write `.bashrc` and `.other`; create unreadable `.private/secret`.
3. Record BackupPath / OtherBackupPath / UnreadableDir.

## Context

- Locks “warn does not abort the job list” — multi-file success after skip.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false)

	const bashBody = "export A=1\n"
	const otherBody = "other-body\n"
	writeFile(t, filepath.Join(home, ".bashrc"), bashBody)
	writeFile(t, filepath.Join(home, ".other"), otherBody)

	private := filepath.Join(home, ".private")
	makeUnreadableDir(t, private, "secret", "nope\n")

	req.Content = bashBody
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.OtherContent = otherBody
	req.OtherSourcePath = filepath.Join(home, ".other")
	req.OtherBackupPath = mappingBackup(req.TargetDir, "alice", ".other")
	req.UnreadableDir = private
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".private", "secret")
	return nil
}
```
