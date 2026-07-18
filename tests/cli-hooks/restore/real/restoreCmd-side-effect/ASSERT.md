# Real restore `restoreCmd` → side-effect file, exit 0

## Expected

- Exit code **0**.
- **SideEffectPath** exists (`restoreCmd` executed).
- Prefer **SourcePath** restored to **Content** when restore still copies the
  file after/alongside restoreCmd (if implementer treats restoreCmd as full
  replacement of copy, marker alone + exit 0 is the hard requirement).

## Side Effects

- `restoreCmd` creates marker via `touch`.
- Backup store content remains readable.

## Errors

- Missing SideEffectPath fails (restoreCmd did not run).
- Non-zero exit fails.

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
	if !pathExists(req.SideEffectPath) {
		t.Fatalf("restoreCmd did not create side-effect %s\nstdout:\n%s\nstderr:\n%s",
			req.SideEffectPath, resp.Stdout, resp.Stderr)
	}
	// Soft preference: if destination was written, it should match backup body.
	if pathExists(req.SourcePath) {
		got := readFileOrEmpty(req.SourcePath)
		if got != req.Content {
			t.Fatalf("restored dest %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
				req.SourcePath, got, req.Content, resp.Stdout, resp.Stderr)
		}
	}
}
```
