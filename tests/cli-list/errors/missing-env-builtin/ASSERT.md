# Missing WORKING_ROLE → non-zero, mentions env

## Expected

- Exit code **non-zero**.
- Combined output mentions **WORKING_ROLE** and ideally “env” / “missing”.

## Side Effects

- No successful full path list of the fixture entries.

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
		t.Fatalf("missing WORKING_ROLE: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(combined, "WORKING_ROLE") {
		t.Fatalf("missing env: output should mention WORKING_ROLE\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
