# Missing validate env W0 → non-zero

## Expected

- Exit code **non-zero**.
- Output mentions **W0** (the missing validate env).

## Side Effects

- None.

## Errors

- Exit 0 fails.

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
		t.Fatalf("missing W0: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(combined, "W0") {
		t.Fatalf("missing validate env: output should mention W0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
