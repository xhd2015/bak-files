# `backup --help` → usage, exit 0

## Expected

- Exit code **0**.
- **Stdout** contains usage text for `backup` (case-insensitive `Usage` and
  the token `backup`).

## Side Effects

- None (no config read; no tree writes).

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
	if !strings.Contains(out, "backup") {
		t.Fatalf("stdout missing backup:\n%s", out)
	}
}
```
