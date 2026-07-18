# basename *.tmp exclude with dots → keep only, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** exists with `keep-me\n`.
- **ExcludedBackupPath** (`noise.tmp`) does **not** exist.

## Side Effects

- Discovered `.myapp` dir is walked; only basename-matched files skipped.

## Errors

- Copying excluded `.tmp` fails the leaf.

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
		t.Fatalf("*.tmp was backed up: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
}
```
