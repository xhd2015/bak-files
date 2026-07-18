# Real backup `beforeCopy` → marker file exists, exit 0

## Expected

- Exit code **0**.
- **MarkerPath** exists (hook executed).
- **BackupPath** equals **Content** (copy still succeeds after hook).

## Side Effects

- `beforeCopy` creates marker via `touch`.
- Source file remains with original content.

## Errors

- Missing marker fails (hook did not run).
- Non-zero exit fails.
- Wrong backup content fails.

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
	if !pathExists(req.MarkerPath) {
		t.Fatalf("beforeCopy did not create marker %s\nstdout:\n%s\nstderr:\n%s",
			req.MarkerPath, resp.Stdout, resp.Stderr)
	}
	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("backup file %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
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
