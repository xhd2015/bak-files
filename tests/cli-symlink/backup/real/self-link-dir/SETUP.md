# Scenario

**Feature**: self-referential symlink-to-dir under a backed-up tree is preserved

```
# operator
HOME=…/home  WORKING_ROLE=alice
mkdir home/Scripts/proj
write home/Scripts/proj/a.txt = "self-link-body\n"
symlink home/Scripts/proj/proj -> abs(home/Scripts/proj)
operator -> bak-files backup --config bak.config.json
  -> exit 0
  -> files/HOME/alice/Scripts/proj/a.txt == body
  -> files/HOME/alice/Scripts/proj/proj is symlink with same target
  -> logs must NOT contain "is a directory"
```

## Preconditions

- Config: `scriptsConfig()` (`~/Scripts`, dots off, mapping role path).
- Self-link points at the **absolute** path of `Scripts/proj` (classic loop case).

## Steps

1. setupScriptsWorld backup real.
2. Create `Scripts/proj/a.txt` and symlink `Scripts/proj/proj` → abs(proj).
3. Record Source/Backup file + link paths and LinkTarget.

## Context

- Primary regression for copySymlink vs Open-follow on dir links.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupScriptsWorld(t, req, "backup", false)

	proj := filepath.Join(home, "Scripts", "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatalf("mkdir proj: %v", err)
	}
	const body = "self-link-body\n"
	srcFile := filepath.Join(proj, "a.txt")
	writeFile(t, srcFile, body)

	// Self-referential: proj/proj -> absolute path of proj directory.
	srcLink := filepath.Join(proj, "proj")
	symlink(t, proj, srcLink)

	req.Content = body
	req.SourceFilePath = srcFile
	req.BackupFilePath = mappingBackup(req.TargetDir, "alice", "Scripts", "proj", "a.txt")
	req.SourceLinkPath = srcLink
	req.BackupLinkPath = mappingBackup(req.TargetDir, "alice", "Scripts", "proj", "proj")
	req.LinkTarget = proj
	return nil
}
```
