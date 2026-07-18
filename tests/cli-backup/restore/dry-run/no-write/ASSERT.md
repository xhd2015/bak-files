# restore --dry-run → exit 0, dest unchanged

## Expected

- Exit code **0**.
- **SourcePath** still absent (or same as DestBefore if it existed).
- Prefer combined output contains `dry-run` or `would`.

## Side Effects

- **Zero writes** to restore destinations.
- Backup store may remain as seeded (read-only).

## Errors

- Creating/modifying dest fails the leaf.
- Non-zero exit fails.

## Exit Code

- **0**

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("dry-run exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	after := readFileOrEmpty(req.SourcePath)
	if after != req.DestBefore {
		t.Fatalf("dry-run modified dest %s:\n got %q\nbefore %q\nstdout:\n%s\nstderr:\n%s",
			req.SourcePath, after, req.DestBefore, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.SourcePath) {
		t.Fatalf("dry-run must not create dest %s", req.SourcePath)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "dry-run") && !strings.Contains(combined, "would") {
		t.Fatalf("dry-run: prefer log with dry-run or would\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
