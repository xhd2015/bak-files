# `"$W0/proj": [".vscode"]` → settings kept, vendor not, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`…/files/W/proj/.vscode/settings.json`) exists with
  **KeepContent** (`vscode-settings\n`).
- **ExcludedBackupPath** (`…/files/W/proj/vendor/poison.go`) does **not** exist.
- No `vendor` directory under the mapped project backup root
  (`…/files/W/proj`).

## Side Effects

- Array whitelist under non-home PREFIX only; never full-PREFIX copyDir.

## Errors

- Missing `.vscode` backup fails.
- Any backed-up vendor path fails (full-PREFIX array-ignored bug).

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
		t.Fatalf("vendor must not be backed up (full-PREFIX array bug): %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
	vendorRoot := filepath.Join(req.TargetDir, "W", "proj", "vendor")
	if pathExists(vendorRoot) {
		t.Fatalf("vendor dir must not appear under store: %s\ntree:\n%s\nstdout:\n%s\nstderr:\n%s",
			vendorRoot, treeFingerprint(t, filepath.Join(req.TargetDir, "W", "proj")),
			resp.Stdout, resp.Stderr)
	}
}
```
