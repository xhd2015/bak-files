# include + exclude same path → skip, exit 0

## Expected

- Exit code **0**.
- **ExcludedBackupPath** does **not** exist under targetDir.

## Side Effects

- Exclude precedence over include for the same path.

## Errors

- Backing up the path fails the leaf.

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
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("exclude must win over include: %s was backed up\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
}
```
