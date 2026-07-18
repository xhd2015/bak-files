# `"~": ["Notes"]` + dots off → Notes only, not Library, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`Notes/a.txt`) exists with **KeepContent**.
- **ExcludedBackupPath** (`Library/x`) does **not** exist.
- No `Library` directory under mapped home store root.

## Side Effects

- Non-dot whitelist entry becomes `~/Notes` job only.

## Errors

- Missing Notes backup or present Library under store fails.

## Exit Code

- **0**

```go
import "testing"

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
		t.Fatalf("Library must not be backed up: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
	libRoot := mappingBackup(req.TargetDir, "alice", "Library")
	if pathExists(libRoot) {
		t.Fatalf("Library dir must not appear under store: %s\nstdout:\n%s\nstderr:\n%s",
			libRoot, resp.Stdout, resp.Stderr)
	}
}
```
