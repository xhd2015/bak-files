# Real backup `cmd` → file under targetDir from stdout, exit 0

## Expected

- Exit code **0**.
- **BackupPath** (`targetDir/HOME/generated.txt`) exists and equals **Content**
  (`hello-from-cmd\n`).

## Side Effects

- Creates `targetDir` tree as needed.
- Command stdout is the sole content of the mapping path file (no leftover
  source-file dependency).

## Errors

- Non-zero exit fails.
- Missing or wrong BackupPath content fails.

## Exit Code

- **0**

```go
import (
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
		t.Fatalf("cmd-generated file %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
}
```
