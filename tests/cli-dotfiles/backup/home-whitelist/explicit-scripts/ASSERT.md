# explicit `~/Scripts` → Scripts kept, Library not, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`Scripts/tool.sh`) exists with **KeepContent**.
- **ExcludedBackupPath** (`Library/x`) does **not** exist.

## Side Effects

- Explicit `~/…` keys remain ordinary jobs; non-listed non-dots stay out.

## Errors

- Missing Scripts copy or Library under store fails.

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
}
```
