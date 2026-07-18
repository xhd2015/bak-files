# `backup --help` → documents --no-dot-files, --include, --exclude

## Expected

- Exit code **0**.
- **Stdout** contains usage text for `backup` and the flag tokens
  `--no-dot-files`, `--include`, and `--exclude` (case-sensitive flag spellings).

## Side Effects

- None (no config read; no tree writes).

## Errors

- Non-zero exit fails.
- Missing any of the three flag tokens fails.

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
	out := resp.Stdout
	if !strings.Contains(strings.ToLower(out), "usage") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	for _, flag := range []string{"--no-dot-files", "--include", "--exclude"} {
		if !strings.Contains(out, flag) {
			t.Fatalf("stdout missing %s:\n%s", flag, out)
		}
	}
}
```
