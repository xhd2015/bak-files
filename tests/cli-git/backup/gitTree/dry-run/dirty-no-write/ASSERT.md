# gitTree --dry-run dirty → no writes, exit 0

## Expected

- Exit code **0**.
- **BackupPath** does not exist.
- TargetDir fingerprint equals **TargetDirBefore**.
- StatsPath content equals **StatsBefore** (empty).
- Combined log preferably contains `dry-run` or `would`.

## Side Effects

- Zero writes under targetDir.
- Zero bak.stats create/update.
- No remote git ops.

## Errors

- Any new file under targetDir or stats change fails.
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
	if pathExists(req.BackupPath) {
		t.Fatalf("dry-run must not write BackupPath %s content=%q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, readFileOrEmpty(req.BackupPath), resp.Stdout, resp.Stderr)
	}
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("dry-run wrote under targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
	stats := readFileOrEmpty(req.StatsPath)
	if stats != req.StatsBefore {
		t.Fatalf("dry-run must not change bak.stats\nbefore %q\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.StatsBefore, stats, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "dry-run") && !strings.Contains(combined, "would") {
		t.Fatalf("dry-run: prefer log with dry-run or would\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
