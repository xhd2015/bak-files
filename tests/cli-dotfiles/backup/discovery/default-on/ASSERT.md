# default-on dry-run → would .bashrc, skip .cache, no writes

## Expected

- Exit code **0**.
- TargetDir fingerprint equals **TargetDirBefore** (still absent / empty).
- **BackupPath** and **ExcludedBackupPath** must not exist as written files.
- Combined stdout/stderr (case-insensitive) mentions `.bashrc` (would/copy intent)
  and a skip signal for `.cache` (prefer `skip` and `cache` / `.cache`).

## Side Effects

- **Zero file writes** under targetDir.

## Errors

- Creating backup files fails the leaf.
- Non-zero exit fails.
- Missing would/skip log tokens fails (encourage implementer logging).

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
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("dry-run wrote under targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.BackupPath) {
		t.Fatalf("dry-run must not create %s", req.BackupPath)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("dry-run must not create skipped path %s", req.ExcludedBackupPath)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, ".bashrc") {
		t.Fatalf("expected log mention of .bashrc\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, "skip") {
		t.Fatalf("expected skip log for pathflag .cache\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, ".cache") && !strings.Contains(combined, "cache") {
		t.Fatalf("expected cache/.cache in skip context\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
