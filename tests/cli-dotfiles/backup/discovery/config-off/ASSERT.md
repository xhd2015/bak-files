# includeDotFiles false → no auto-dot backup, exit 0

## Expected

- Exit code **0**.
- **BackupPath** (`…/HOME/alice/.bashrc`) does **not** exist.

## Side Effects

- No discovered-dot copies under targetDir for home top-level dots.

## Errors

- Copying `.bashrc` via auto-discovery fails the leaf.

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
		t.Fatalf("auto-dot must not backup when includeDotFiles false: %s\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, resp.Stdout, resp.Stderr)
	}
}
```
