# Real restore → dest file content, exit 0

## Expected

- Exit code **0**.
- **SourcePath** (restore destination) exists with content `from-backup\n`.
- Seeded **BackupPath** still holds the same content (restore is copy, not move).

## Side Effects

- Creates parent dirs for dest as needed.
- Must not delete the backup store file.

## Errors

- Non-zero exit fails.
- Wrong or missing dest content fails.

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
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	got := readFileOrEmpty(req.SourcePath)
	if got != req.Content {
		t.Fatalf("restored dest %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.SourcePath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	bak, err := os.ReadFile(req.BackupPath)
	if err != nil {
		t.Fatalf("backup store missing after restore: %v", err)
	}
	if string(bak) != req.Content {
		t.Fatalf("backup store altered: got %q want %q", bak, req.Content)
	}
}
```
