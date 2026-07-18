# Scenario

**Feature**: restore dry-run leaves destination unchanged

```
# operator
seed files/HOME/notes.txt = "from-backup\n"
dest home/notes.txt absent
operator -> bak-files restore --config … --dry-run
  -> exit 0
  -> dest still absent
  -> prefer log dry-run or would
```

## Preconditions

- Backup seeded; destination absent before Run.

## Steps

1. setupSimpleBackupWorld restore + dry-run.
2. Seed BackupPath; remove SourcePath.
3. DestBefore empty string (absent).

## Context

- Core P3 dry-run exit criterion for restore.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	const body = "from-backup\n"
	setupSimpleBackupWorld(t, req, "restore", true, body)
	writeFile(t, req.BackupPath, body)
	_ = os.Remove(req.SourcePath)
	req.Content = body
	req.DestBefore = readFileOrEmpty(req.SourcePath)
	return nil
}
```
