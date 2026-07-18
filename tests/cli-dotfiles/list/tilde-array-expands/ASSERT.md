# list expands `"~": [".bashrc"]` → HOME/alice/.bashrc, not bare HOME/alice

## Expected

- Exit code **0**.
- **Stdout** contains an exact trimmed line equal to **ExpectedMappingPath**
  (`HOME/alice/.bashrc`).
- **Stdout** must **not** contain an exact trimmed line `HOME/alice` (that is the
  mapping path of an unexpanded bare `"~"` full-home job).

## Side Effects

- No targetDir requirement.

## Errors

- Missing expanded path or presence of bare home mapping line fails.

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
	foundExpanded := false
	foundBareHome := false
	for _, line := range strings.Split(resp.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == req.ExpectedMappingPath {
			foundExpanded = true
		}
		if line == "HOME/alice" {
			foundBareHome = true
		}
	}
	if !foundExpanded {
		t.Fatalf("stdout missing expanded mapping path %q\nstdout:\n%s\nstderr:\n%s",
			req.ExpectedMappingPath, resp.Stdout, resp.Stderr)
	}
	if foundBareHome {
		t.Fatalf("stdout must not list bare home mapping HOME/alice (unexpanded \"~\" job)\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
