# Missing config file → non-zero

## Expected

- Exit code **non-zero**.
- **Stderr** (or combined output) indicates the config/file problem (e.g.
  mentions `config`, `no such`, `not found`, `open`, or the path basename).

## Side Effects

- No mapping-path listing treated as success (stdout should not look like a
  clean multi-entry path list; empty or error-only is fine).

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
		t.Fatalf("missing config: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	ok := strings.Contains(combined, "config") ||
		strings.Contains(combined, "not found") ||
		strings.Contains(combined, "no such") ||
		strings.Contains(combined, "open") ||
		strings.Contains(combined, "missing") ||
		strings.Contains(combined, "exist")
	if !ok {
		t.Fatalf("missing config: stderr/stdout should mention file/config error\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
