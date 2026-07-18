# Failing `beforeCopy` → non-zero exit

## Expected

- Exit code **non-zero** (hook failure fails the whole run).
- Prefer not treating the run as success; BackupPath need not exist (if
  beforeCopy aborts before copy, backup file should be absent or unchanged).

## Side Effects

- Process must not exit 0 after a failing beforeCopy.

## Errors

- Exit 0 fails this leaf (silent ignore of hook failure is a bug).

## Exit Code

- **≠ 0**

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		// exec.ExitError is converted to ExitCode by Run; other errors are hard failures.
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("failing beforeCopy must exit non-zero, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
