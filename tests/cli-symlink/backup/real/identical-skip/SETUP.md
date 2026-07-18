# Scenario

**Feature**: second backup of identical symlink tree exits 0 and leaves links stable

```
# operator
# first backup already produced targetDir tree with same file + symlink
operator -> bak-files backup --config …   # second run
  -> exit 0
  -> backup symlink target still LinkTarget
  -> a.txt body unchanged
  (optional: verbose identical symlink skip; not required)
```

## Preconditions

- Self-link-dir style fixture (file + abs self-dir link under `Scripts/proj`).
- First real backup runs **in Setup** via the same binary harness pattern
  (inline `go run` / pre-seed via first exec through ensure path).
- Prefer: pre-seed by calling the built binary once in Setup after fixtures.

## Steps

1. Build self-link-dir world (real backup Args).
2. Run first backup in Setup (must succeed) to populate targetDir.
3. Snapshot optional TargetDirBefore fingerprint (symlink-aware).
4. Leaf Run executes the second backup with the same Args.

## Context

- Coverage for `sameContent` symlink branch + identical symlink early return.
- Do not invent a must-fail assertion; GREEN if second run is stable.

```go
import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupScriptsWorld(t, req, "backup", false)

	proj := filepath.Join(home, "Scripts", "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatalf("mkdir proj: %v", err)
	}
	const body = "identical-symlink-body\n"
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

	// First backup: populate targetDir so leaf Run is a second pass.
	bin := ensureBakFilesBinary(t)
	cmd := exec.Command(bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = req.Env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("first backup failed: %v\nstdout:\n%s\nstderr:\n%s",
			err, stdout.String(), stderr.String())
	}
	if !isSymlink(req.BackupLinkPath) {
		t.Fatalf("first backup did not create symlink at %s\nstdout:\n%s\nstderr:\n%s",
			req.BackupLinkPath, stdout.String(), stderr.String())
	}
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	return nil
}
```
