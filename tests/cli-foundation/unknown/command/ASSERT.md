# unknown command → non-zero, stderr error prefix

## Expected

- Exit code **non-zero**.
- **Stderr** contains `Error:` or `bak-files:` (error surfaced on stderr).

## Side Effects

- None.

## Errors

- Exit 0 or silent stderr fails the contract.

## Exit Code

- **≠ 0**

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("unknown command: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	errOut := resp.Stderr
	if !strings.Contains(errOut, "Error:") && !strings.Contains(errOut, "bak-files:") {
		// also accept program name prefix variants without colon spacing
		if !containsFold(errOut, "error") && !strings.Contains(errOut, "bak-files") {
			t.Fatalf("unknown command: stderr missing Error: or bak-files: prefix\nstderr:\n%s\nstdout:\n%s",
				resp.Stderr, resp.Stdout)
		}
	}
}
```
