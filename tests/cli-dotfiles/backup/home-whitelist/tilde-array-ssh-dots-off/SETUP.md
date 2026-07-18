# Scenario

**Feature**: dots off + `"~": [".ssh"]` expands only to `~/.ssh`, not other dots

```
# operator
includeDotFiles: false
files: { "~": [".ssh"] }
HOME: .ssh/config and .bashrc both present
operator -> bak-files backup --config …
  -> files/HOME/alice/.ssh/config == "Host *\n"
  -> files/HOME/alice/.bashrc MUST NOT exist
  -> exit 0
```

## Preconditions

- `tildeArrayConfigDotsOff(['.ssh'])`.
- No synthetic `~/.*` discovery when dots off.

## Steps

1. setupDotsWorld tildeArrayConfigDotsOff, real backup.
2. Plant `.ssh/config` and `.bashrc`.
3. Keep = `.ssh/config`; Excluded = `.bashrc` under store.

## Context

- Proves `"~"` array names become ordinary `$HOME/name` jobs when listed.
- `.bashrc` must not ride along via full-home walk or auto-discovery.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, tildeArrayConfigDotsOff(`[".ssh"]`), "backup", false)
	const sshBody = "Host *\n"
	writeFile(t, filepath.Join(home, ".ssh", "config"), sshBody)
	writeFile(t, filepath.Join(home, ".bashrc"), "should-not-backup\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".ssh", "config")
	req.KeepContent = sshBody
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	return nil
}
```
