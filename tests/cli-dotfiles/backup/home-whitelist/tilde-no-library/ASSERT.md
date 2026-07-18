# `"~": [".bashrc"]` → bashrc kept, Library not, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`…/HOME/alice/.bashrc`) exists with **KeepContent** (`rc-body\n`).
- **ExcludedBackupPath** (`…/HOME/alice/Library/Android/sdk/x`) does **not** exist.
- No `Library` tree under the mapped home backup root.

## Side Effects

- Whitelist + optional synthetic `~/.*` discovery only; never full-home copyDir.

## Errors

- Missing bashrc backup fails.
- Any backed-up Library path fails (full-home bug).

## Exit Code

- **0**

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	got := readFileOrEmpty(req.KeepBackupPath)
	if got != req.KeepContent {
		t.Fatalf("keep %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.KeepBackupPath, got, req.KeepContent, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("Library must not be backed up (full-home bug): %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
	libRoot := mappingBackup(req.TargetDir, "alice", "Library")
	if pathExists(libRoot) {
		t.Fatalf("Library dir must not appear under store: %s\ntree:\n%s\nstdout:\n%s\nstderr:\n%s",
			libRoot, treeFingerprint(t, filepath.Join(req.TargetDir, "HOME", "alice")),
			resp.Stdout, resp.Stderr)
	}
}
```
