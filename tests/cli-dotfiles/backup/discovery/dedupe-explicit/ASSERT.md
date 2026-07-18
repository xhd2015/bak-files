# explicit + discover → one copy under mapping path, exit 0

## Expected

- Exit code **0**.
- **BackupPath** exists with content **Content** (`once\n`).
- Source still present with same content (copy, not move).

## Side Effects

- Exactly the expected mapping path is populated; no requirement to count jobs,
  but content must match a single successful write of the fixture body.

## Errors

- Missing/wrong BackupPath content fails.

## Exit Code

- **0**

```go
import (
	"os"
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
	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	src, err := os.ReadFile(req.SourcePath)
	if err != nil {
		t.Fatalf("source missing after backup: %v", err)
	}
	if string(src) != req.Content {
		t.Fatalf("source altered: got %q want %q", src, req.Content)
	}
}
```
