# list discovers `.bashrc` → mapping path on stdout, exit 0

## Expected

- Exit code **0**.
- **Stdout** contains a line with **ExpectedMappingPath** (`HOME/alice/.bashrc`).

## Expected Output

Stdout includes (among optional other discovered paths) a line:

```
HOME/alice/.bashrc
```

Trailing newline preferred for CLI list output.

## Side Effects

- No requirement to create targetDir.

## Errors

- Non-zero exit fails.
- Missing ExpectedMappingPath fails.

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
	found := false
	for _, line := range strings.Split(resp.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == req.ExpectedMappingPath {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("stdout missing mapping path %q\nstdout:\n%s\nstderr:\n%s",
			req.ExpectedMappingPath, resp.Stdout, resp.Stderr)
	}
}
```
