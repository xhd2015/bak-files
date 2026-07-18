# bare `list` without config → Error: open config

## Expected

- Exit code **non-zero**.
- **Stderr** contains `Error:` (error prefix).
- **Stderr** mentions config open failure: substring `config` and `bak.config.json`
  (e.g. `Error: open config bak.config.json: …`).

## Side Effects

- None (no successful list output; no config written).

## Errors

- Exit 0 would mean list succeeded without a config (incorrect for this cwd).
- "not implemented" is **not** expected — list is real as of P2.

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
		t.Fatalf("list without config: want non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	errOut := resp.Stderr
	if !strings.Contains(errOut, "Error:") {
		t.Fatalf("list without config: stderr missing Error: prefix\nstderr:\n%s\nstdout:\n%s",
			resp.Stderr, resp.Stdout)
	}
	if !containsFold(errOut, "config") {
		t.Fatalf("list without config: stderr missing config mention\nstderr:\n%s", errOut)
	}
	if !strings.Contains(errOut, "bak.config.json") {
		t.Fatalf("list without config: stderr missing bak.config.json\nstderr:\n%s", errOut)
	}
}
```
