# Scenario

**Feature**: `--include .cache` force-keeps a pathflag-cache tree

```
# operator
files={}; write home/.cache/keep.txt = "cached\n"
operator -> bak-files backup --config … --include .cache
  -> files/HOME/alice/.cache/keep.txt == "cached\n"
  -> exit 0
```

## Preconditions

- Without include, `.cache` would be pathflag-skipped.

## Steps

1. setupDotsWorld emptyFilesConfig, backup, real, `--include`, `.cache`.
2. Write `.cache/keep.txt`; set KeepBackupPath / Content.

## Context

- force-include is step 1 of skip policy (before pathflag).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false,
		"--include", ".cache")
	const body = "cached\n"
	writeFile(t, filepath.Join(home, ".cache", "keep.txt"), body)
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".cache", "keep.txt")
	req.KeepContent = body
	req.Content = body
	return nil
}
```
