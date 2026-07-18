# pathflag partial → config kept, .tmp skipped, exit 0

## Expected

- Exit code **0**.
- **KeepBackupPath** (`.codex/config`) exists with `cfg\n`.
- **ExcludedBackupPath** (`.codex/.tmp/plugin`) does **not** exist.
- Prefer combined logs contain `skip` and a path mentioning `.tmp` or `codex`
  (soft preference strengthened to fail if no skip token at all).

## Side Effects

- Directory job `.codex` walks; only flagged subpaths skipped.

## Errors

- Missing config or present tmp under target fails.

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
	got := readFileOrEmpty(req.KeepBackupPath)
	if got != req.KeepContent {
		t.Fatalf("keep %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.KeepBackupPath, got, req.KeepContent, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("pathflag tmp was backed up: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "skip") {
		t.Fatalf("prefer skip log for pathflag tmp path\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
