# list expands `"$W0/proj": [".vscode"]` → W/proj/.vscode, not bare W/proj

## Expected

- Exit code **0**.
- **Stdout** contains an exact trimmed line equal to **ExpectedMappingPath**
  (`W/proj/.vscode`).
- **Stdout** must **not** contain an exact trimmed line `W/proj` (that is the
  mapping path of an unexpanded full-PREFIX job for `"$W0/proj"`).

## Side Effects

- No targetDir requirement.

## Errors

- Missing expanded child path or presence of bare PREFIX mapping line fails.

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
	foundBarePrefix := false
	for _, line := range strings.Split(resp.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == req.ExpectedMappingPath {
			foundExpanded = true
		}
		if line == "W/proj" {
			foundBarePrefix = true
		}
	}
	if !foundExpanded {
		t.Fatalf("stdout missing expanded mapping path %q\nstdout:\n%s\nstderr:\n%s",
			req.ExpectedMappingPath, resp.Stdout, resp.Stderr)
	}
	if foundBarePrefix {
		t.Fatalf("stdout must not list bare PREFIX mapping W/proj (unexpanded array job)\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
