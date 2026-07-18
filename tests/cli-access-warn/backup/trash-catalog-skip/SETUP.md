# Scenario

**Feature**: pathflag `.Trash` is not scheduled as a job; bashrc still backed up

```
# operator
HOME=…/home  files={}
mkdir home/.Trash          # empty dir enough; pathflag rule
write home/.bashrc = "export TRASH_SKIP=1\n"
operator -> bak-files backup --config …
  -> exit 0
  -> log INFO: skip … .Trash … (macOS trash)  [or skip + trash reason]
  -> files/HOME/alice/.bashrc == body
  -> no .Trash tree under targetDir
```

## Preconditions

- Empty files; includeDotFiles omitted (default true).
- `.Trash` is pathflag Trash → DefaultSkipMask; discovery logSkip, no Job.
- Empty `.Trash` directory is enough to appear in ReadDir.

## Steps

1. setupDotsWorld emptyFilesConfig, backup, dryRun=false.
2. Create empty `home/.Trash` and write `.bashrc`.
3. Record BackupPath for bashrc; ExcludedBackupPath for any `.Trash` under mapping.

## Context

- Catalog skip at discovery (preferred over walk open-fail on real macOS trash).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false)

	trash := filepath.Join(home, ".Trash")
	if err := os.MkdirAll(trash, 0o755); err != nil {
		t.Fatalf("mkdir .Trash: %v", err)
	}

	const body = "export TRASH_SKIP=1\n"
	writeFile(t, filepath.Join(home, ".bashrc"), body)

	req.Content = body
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.BackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	// Any path under mapping for .Trash must stay absent.
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".Trash")
	return nil
}
```
