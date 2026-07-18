# restore dry-run pathflag → no dest write; prefer skip log for .tmp

## Expected

- Exit code **0**.
- **DestPath** (`home/.codex/.tmp/plugin`) remains absent / DestBefore.
- **SourcePath** (config dest) also remains absent (dry-run).
- Combined logs prefer `skip` (pathflag) for the tmp path.

## Side Effects

- **Zero writes** to restore destinations under HOME.

## Errors

- Creating DestPath or SourcePath fails.
- Non-zero exit fails.
- Missing skip token fails (encourage always-on skip INFO).

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
		t.Fatalf("exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.DestPath) {
		t.Fatalf("dry-run restore must not write tmp dest %s\nstdout:\n%s\nstderr:\n%s",
			req.DestPath, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.SourcePath) {
		t.Fatalf("dry-run restore must not write config dest %s\nstdout:\n%s\nstderr:\n%s",
			req.SourcePath, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "skip") {
		t.Fatalf("prefer skip log for pathflag tmp on restore dry-run\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, ".tmp") && !strings.Contains(combined, "tmp") {
		t.Fatalf("prefer tmp path in skip context\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
