# Invalid JSON → non-zero

## Expected

- Exit code **non-zero**.
- Output mentions parse/config/JSON/invalid (case-insensitive).

## Side Effects

- None required.

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
		t.Fatalf("invalid JSON: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	ok := strings.Contains(combined, "parse") ||
		strings.Contains(combined, "json") ||
		strings.Contains(combined, "config") ||
		strings.Contains(combined, "invalid") ||
		strings.Contains(combined, "syntax")
	if !ok {
		t.Fatalf("invalid JSON: expected parse/config error signal\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
