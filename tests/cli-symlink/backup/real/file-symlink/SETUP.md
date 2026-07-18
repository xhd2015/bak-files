# Scenario

**Feature**: file symlink is recreated as a symlink under the mapping path

```
# operator
write home/Scripts/target.txt = "link-target-body\n"
symlink home/Scripts/link.txt -> "target.txt"   # relative target string
operator -> bak-files backup --config …
  -> files/HOME/alice/Scripts/target.txt == body
  -> files/HOME/alice/Scripts/link.txt is symlink -> "target.txt"
  -> exit 0
```

## Preconditions

- Relative link target string must be preserved (not necessarily rewritten).
- Real backup (no dry-run).

## Steps

1. setupScriptsWorld backup real.
2. Write target.txt; create relative symlink link.txt.
3. Record paths and LinkTarget `"target.txt"`.

## Context

- Asserts preserve-as-link, not Open-and-copy through the link as the only mode
  for the link path itself.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupScriptsWorld(t, req, "backup", false)

	const body = "link-target-body\n"
	srcFile := filepath.Join(home, "Scripts", "target.txt")
	writeFile(t, srcFile, body)

	srcLink := filepath.Join(home, "Scripts", "link.txt")
	const relTarget = "target.txt"
	symlink(t, relTarget, srcLink)

	req.Content = body
	req.SourceFilePath = srcFile
	req.BackupFilePath = mappingBackup(req.TargetDir, "alice", "Scripts", "target.txt")
	req.SourceLinkPath = srcLink
	req.BackupLinkPath = mappingBackup(req.TargetDir, "alice", "Scripts", "link.txt")
	req.LinkTarget = relTarget
	return nil
}
```
