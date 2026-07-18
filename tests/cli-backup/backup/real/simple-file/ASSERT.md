# Real backup → file under targetDir, exit 0

## Expected

- Exit code **0**.
- **BackupPath** (`targetDir/HOME/notes.txt`) exists and equals **Content**
  (`hello-backup\n`).
- Source file remains readable with the same content (backup is copy, not move).

## Side Effects

- Creates `targetDir` tree as needed.
- Must not leave source deleted.

## Errors

- Non-zero exit fails.
- Missing or wrong BackupPath content fails.

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
