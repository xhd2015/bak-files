# empty argv → root Usage

## Expected

- Exit code **0**.
- **Stdout** mentions usage (case-insensitive `usage`) and the three subcommands
  `backup`, `restore`, and `list`.

## Side Effects

- None (no files written).

## Errors

- `Run` harness error is fatal.
- Non-zero exit or missing usage tokens fail the test.

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
	if !containsFold(out, "usage") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	for _, cmd := range []string{"backup", "restore", "list"} {
		if !strings.Contains(out, cmd) {
			t.Fatalf("stdout missing subcommand %q:\n%s", cmd, out)
		}
	}
}
```
