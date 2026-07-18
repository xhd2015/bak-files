# Backup with excludes → keep.txt only, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** exists with content `keep-me\n`.
- **ExcludedBackupPath** (`…/noise.tmp`) does **not** exist under targetDir.

## Side Effects

- Directory may exist under targetDir; excluded basename must be absent.

## Errors

- Copying the excluded file fails the leaf.
- Non-zero exit fails.

## Exit Code

- **0**

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	got := readFileOrEmpty(req.KeepBackupPath)
	if got != req.Content {
		t.Fatalf("kept backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.KeepBackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("excluded file was backed up: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
}
```
