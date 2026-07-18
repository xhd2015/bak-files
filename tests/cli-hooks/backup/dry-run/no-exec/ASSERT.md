# backup --dry-run → no hook/`cmd` exec, no writes, exit 0

## Expected

- Exit code **0**.
- **MarkerPath** does not exist (`beforeCopy` did not run).
- **BackupPath** (cmd-generated file) does not exist.
- TargetDir fingerprint equals **TargetDirBefore** (no new/changed files).
- Prefer combined output contains `dry-run` or `would`.

## Side Effects

- Zero execution of shell hooks/cmds for these entries.
- Zero file writes under targetDir for generated/copied content.

## Errors

- Creating marker or BackupPath fails the leaf.
- Non-zero exit fails.
- Silent success without dry-run/would log fails (encourage logging).

## Exit Code

- **0**

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("dry-run exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.MarkerPath) {
		t.Fatalf("dry-run must not run beforeCopy; marker exists: %s\nstdout:\n%s\nstderr:\n%s",
			req.MarkerPath, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.BackupPath) {
		t.Fatalf("dry-run must not run cmd; file exists: %s content=%q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, readFileOrEmpty(req.BackupPath), resp.Stdout, resp.Stderr)
	}
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("dry-run wrote under targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
	// Also ensure the plain-file mapping path was not written.
	notesBackup := filepath.Join(req.TargetDir, "HOME", "notes.txt")
	if pathExists(notesBackup) {
		t.Fatalf("dry-run must not copy notes to %s", notesBackup)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "dry-run") && !strings.Contains(combined, "would") {
		t.Fatalf("dry-run: prefer log with dry-run or would\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
