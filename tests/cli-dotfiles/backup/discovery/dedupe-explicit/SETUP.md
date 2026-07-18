# Scenario

**Feature**: explicit job wins; discovery does not double-copy same source

```
# operator
files: {"~/.bashrc": true}; home/.bashrc = "once\n"
operator -> bak-files backup --config …
  -> files/HOME/alice/.bashrc == "once\n"
  -> exit 0 (single logical copy)
```

## Preconditions

- explicitBashrcConfig; source present.

## Steps

1. setupDotsWorld with explicitBashrcConfig, real backup.
2. Write home/.bashrc; set BackupPath and Content.

## Context

- Locked decision: explicit job wins over auto-discover for same source.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, explicitBashrcConfig(), "backup", false)
	const body = "once\n"
	writeFile(t, filepath.Join(home, ".bashrc"), body)
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.Content = body
	return nil
}
```
