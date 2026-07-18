# Restore after gitTree backup seed → source content, exit 0

## Expected

- Exit code **0**.
- **SourcePath** exists and equals **Content** (`from-backup\n`).
- **BackupPath** still readable with the same content (restore is copy, not move).

## Side Effects

- Creates/overwrites operator-side README under `~/repo`.
- Does not require network remotes or git FF.

## Errors

- Non-zero exit fails.
- Wrong or missing SourcePath content fails.

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
		t.Fatalf("restored source %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.SourcePath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	bak, err := os.ReadFile(req.BackupPath)
	if err != nil {
		t.Fatalf("backup missing after restore: %v", err)
	}
	if string(bak) != req.Content {
		t.Fatalf("backup altered by restore: got %q want %q", bak, req.Content)
	}
}
```
