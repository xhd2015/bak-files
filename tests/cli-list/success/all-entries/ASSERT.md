# Full list → two mapping paths, exit 0

## Expected

- Exit code **0**.
- **Stdout** is exactly the two mapping paths in fixture `files` order,
  one per line, trailing newline (POSIX CLI).

## Expected Output

```
HOME/alice/.bashrc
HOME/Scripts/tool.sh
```

## Side Effects

- No requirement to create `targetDir` or copy files for list.

## Errors

- Non-zero exit fails.
- Extra/missing lines fail the template match.

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
