# Scenario

**Feature**: `global.includeDotFiles: false` disables auto-dot discovery

```
# operator
config includeDotFiles: false; files={}
write home/.bashrc = "export B=1\n"
operator -> bak-files backup --config …
  -> exit 0
  -> files/HOME/alice/.bashrc MUST NOT exist
```

## Preconditions

- emptyFilesConfigDotsOff; home has `.bashrc`.

## Steps

1. setupDotsWorld with dots-off config, real backup.
2. Write `.bashrc`; record BackupPath.

## Context

- Config-off is independent of CLI `--no-dot-files`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfigDotsOff(), "backup", false)
	writeFile(t, filepath.Join(home, ".bashrc"), "export B=1\n")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.Content = "export B=1\n"
	return nil
}
```
