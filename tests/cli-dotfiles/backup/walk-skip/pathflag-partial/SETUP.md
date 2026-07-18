# Scenario

**Feature**: pathflag partial — keep `.codex/config`, skip `.codex/.tmp`

```
# operator
write home/.codex/config = "cfg\n"
write home/.codex/.tmp/plugin = "tmp\n"
files={}  # discover ~/.codex as dir job
operator -> bak-files backup --config …
  -> files/HOME/alice/.codex/config == "cfg\n"
  -> .codex/.tmp/plugin MUST NOT exist under targetDir
  -> prefer INFO skip log for .codex/.tmp
  -> exit 0
```

## Preconditions

- `.codex/.tmp` is pathflag Tmp|Cache (DefaultSkipMask); config is owner-only/no skip.

## Steps

1. setupDotsWorld emptyFilesConfig, real backup.
2. Plant config and .tmp/plugin; record Keep and Excluded paths.

## Context

- Demonstrates walk-time pathflag, not whole-job drop of `.codex`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false)
	writeFile(t, filepath.Join(home, ".codex", "config"), "cfg\n")
	writeFile(t, filepath.Join(home, ".codex", ".tmp", "plugin"), "tmp\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".codex", "config")
	req.KeepContent = "cfg\n"
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".codex", ".tmp", "plugin")
	return nil
}
```
