# Scenario

**Feature**: backup dry-run leaves targetDir untouched

```
# operator
write home/notes.txt = "dry-run-src\n"
targetDir absent
operator -> bak-files backup --config … --dry-run
  -> exit 0
  -> targetDir still absent OR fingerprint unchanged
  -> prefer stdout/stderr contains "dry-run" or "would"
```

## Preconditions

- Simple file fixture; source present; no pre-existing targetDir tree.

## Steps

1. setupSimpleBackupWorld with dryRun=true, content `dry-run-src\n`.
2. TargetDirBefore = "" (absent).

## Context

- Core P3 dry-run exit criterion for backup.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	setupSimpleBackupWorld(t, req, "backup", true, "dry-run-src\n")
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	return nil
}
```
