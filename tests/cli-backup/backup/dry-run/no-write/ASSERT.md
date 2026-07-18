# backup --dry-run → exit 0, no targetDir writes

## Expected

- Exit code **0**.
- TargetDir fingerprint equals **TargetDirBefore** (still absent / empty).
- **BackupPath** must not exist as a written copy of Content.
- Prefer combined output contains `dry-run` or `would` (soft if tree is clean).

## Side Effects

- **Zero file writes** under targetDir (no new files, no overwrites).

## Errors

- Creating BackupPath fails the leaf.
- Non-zero exit fails.

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
		t.Fatalf("dry-run exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("dry-run wrote under targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.BackupPath) {
		t.Fatalf("dry-run must not create %s", req.BackupPath)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "dry-run") && !strings.Contains(combined, "would") {
		// Soft preference: still pass if tree clean, but fail to encourage logging.
		t.Fatalf("dry-run: prefer log with dry-run or would\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
