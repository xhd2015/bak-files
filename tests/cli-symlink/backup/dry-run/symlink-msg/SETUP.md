# Scenario

**Feature**: dry-run on self-link tree mentions would symlink and writes nothing

```
# operator
# self-link-dir fixture under ~/Scripts/proj
operator -> bak-files backup --config … --dry-run
  -> exit 0
  -> targetDir absent / fingerprint unchanged
  -> combined log contains "would symlink" (or dry-run + symlink)
```

## Preconditions

- Same self-link-dir source as real leaf; no pre-existing targetDir.
- Args include `--dry-run`.

## Steps

1. setupScriptsWorld backup dryRun=true.
2. Create proj/a.txt + self symlink proj/proj.
3. TargetDirBefore = treeFingerprint (expect empty).

## Context

- Engine log: `dry-run: would symlink %s -> %s`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupScriptsWorld(t, req, "backup", true)

	proj := filepath.Join(home, "Scripts", "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatalf("mkdir proj: %v", err)
	}
	const body = "dry-run-symlink-body\n"
	srcFile := filepath.Join(proj, "a.txt")
	writeFile(t, srcFile, body)

	srcLink := filepath.Join(proj, "proj")
	symlink(t, proj, srcLink)

	req.Content = body
	req.SourceFilePath = srcFile
	req.BackupFilePath = mappingBackup(req.TargetDir, "alice", "Scripts", "proj", "a.txt")
	req.SourceLinkPath = srcLink
	req.BackupLinkPath = mappingBackup(req.TargetDir, "alice", "Scripts", "proj", "proj")
	req.LinkTarget = proj
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	return nil
}
```
