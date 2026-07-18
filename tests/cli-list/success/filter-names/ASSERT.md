# NAMES filter → single mapping path, exit 0

## Expected

- Exit code **0**.
- **Stdout** is only `HOME/Scripts/tool.sh` (one line + trailing newline).

## Expected Output

```
HOME/Scripts/tool.sh
```

## Side Effects

- None.

## Errors

- Including the filtered-out `.bashrc` path fails.

## Exit Code

- **0**

```go
import (
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
HOME/Scripts/tool\.sh
`)
}
```
