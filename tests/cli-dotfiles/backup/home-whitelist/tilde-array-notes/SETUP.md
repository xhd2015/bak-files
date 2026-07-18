# Scenario

**Feature**: non-dot names in `"~"` array are whitelist additions (`~/Notes`)

```
# operator
includeDotFiles: false
files: { "~": ["Notes"] }
HOME: Notes/a.txt and Library/x
operator -> bak-files backup --config …
  -> files/HOME/alice/Notes/a.txt == "note-a\n"
  -> Library MUST NOT be under store
  -> exit 0
```

## Preconditions

- Dots off so only the expanded `Notes` job is scheduled from `"~"`.
- Poison `Library/x` proves no full-home Source=$HOME.

## Steps

1. setupDotsWorld tildeArrayConfigDotsOff `["Notes"]`, real backup.
2. Plant Notes and Library; record keep/excluded paths.

## Context

- Non-dot basename under `"~"` → ordinary job `$HOME/Notes`, mapping `HOME/alice/Notes`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, tildeArrayConfigDotsOff(`["Notes"]`), "backup", false)
	const body = "note-a\n"
	writeFile(t, filepath.Join(home, "Notes", "a.txt"), body)
	writeFile(t, filepath.Join(home, "Library", "x"), "lib-poison\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", "Notes", "a.txt")
	req.KeepContent = body
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", "Library", "x")
	return nil
}
```
