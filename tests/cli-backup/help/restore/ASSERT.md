# `restore --help` → usage, exit 0

## Expected

- Exit code **0**.
- **Stdout** contains usage text for `restore` (case-insensitive `Usage` and
  the token `restore`).

## Side Effects

- None.

## Errors

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
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout
	if !strings.Contains(strings.ToLower(out), "usage") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	if !strings.Contains(out, "restore") {
		t.Fatalf("stdout missing restore:\n%s", out)
	}
}
```
