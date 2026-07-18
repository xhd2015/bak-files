# `list --help` → usage, exit 0

## Expected

- Exit code **0**.
- **Stdout** contains usage text for the `list` command (case-insensitive
  `Usage` and the token `list`).
- Preferably mentions `--config` (soft: if missing, still pass once Usage+list
  present — implementer should document `--config` for P2).

## Expected Output

Help wording may evolve; assert via structured checks rather than a full
template of the entire help page.

## Side Effects

- None (no config read required).

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
	if !strings.Contains(out, "list") {
		t.Fatalf("stdout missing list:\n%s", out)
	}
}
```
