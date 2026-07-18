# Scenario

**Feature**: restore writes destination file from seeded targetDir

```
# operator
seed files/HOME/notes.txt = "from-backup\n"
dest home/notes.txt absent
operator -> bak-files restore --config …
  -> home/notes.txt == "from-backup\n"
  -> exit 0
```

## Preconditions

- Config `simpleFileConfig()`; backup store pre-seeded; source dest absent.

## Steps

1. setupSimpleBackupWorld for restore (do not create source; empty content flag).
2. Write BackupPath with `from-backup\n`.
3. Ensure SourcePath does not exist.

## Context

- Primary exit criterion for P3 restore happy path.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	const body = "from-backup\n"
	// content "" skips writing source in helper when command is backup only —
	// helper writes source only for command=="backup". For restore, seed backup.
	setupSimpleBackupWorld(t, req, "restore", false, body)
	writeFile(t, req.BackupPath, body)
	// Ensure destination is absent before restore.
	_ = os.Remove(req.SourcePath)
	req.Content = body
	return nil
}
```
