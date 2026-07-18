# Scenario

**Feature**: `--no-dot-files` disables auto-dot discovery on the CLI

```
# operator
files={}; home/.bashrc present; default includeDotFiles true
operator -> bak-files backup --config … --no-dot-files
  -> exit 0
  -> files/HOME/alice/.bashrc MUST NOT exist
```

## Preconditions

- emptyFilesConfig (dots would be on without the flag).

## Steps

1. setupDotsWorld emptyFilesConfig, backup, dryRun=false, extra `--no-dot-files`.
2. Write `.bashrc`; record BackupPath.

## Context

- CLI disable is **`--no-dot-files` only** (no `--dot-files` enable flag).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false, "--no-dot-files")
	writeFile(t, filepath.Join(home, ".bashrc"), "export C=1\n")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.Content = "export C=1\n"
	return nil
}
```
