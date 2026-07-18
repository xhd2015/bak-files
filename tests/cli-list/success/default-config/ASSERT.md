# Default bak.config.json → full path list, exit 0

## Expected

- Exit code **0**.
- Same stdout as `success/all-entries` (fixture order mapping paths).

## Expected Output

```
HOME/alice/.bashrc
HOME/Scripts/tool.sh
```

## Side Effects

- None.

## Errors

- Non-zero exit fails.

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
HOME/alice/\.bashrc
HOME/Scripts/tool\.sh
`)
}
```
