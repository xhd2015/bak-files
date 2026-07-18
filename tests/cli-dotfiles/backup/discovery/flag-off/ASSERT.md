# --no-dot-files → no auto-dot backup, exit 0

## Expected

- Exit code **0**.
- **BackupPath** does **not** exist.

## Side Effects

- Flag overrides default-on discovery for this invocation.

## Errors

- Auto-backing up `.bashrc` fails the leaf.

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
	if pathExists(req.BackupPath) {
		t.Fatalf("--no-dot-files must not auto-backup: %s\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, resp.Stdout, resp.Stderr)
	}
}
```
