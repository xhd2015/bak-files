# backup --dry-run → exit 0, would symlink, no targetDir writes

## Expected

- Exit code **0**.
- TargetDir fingerprint equals **TargetDirBefore** (still absent / empty).
- **BackupLinkPath** and **BackupFilePath** must not exist as created copies.
- Combined output contains **`would symlink`** (preferred exact engine token)
  or both `dry-run` and `symlink` as a softer fallback — prefer failing if
  neither form of messaging is present.

## Side Effects

- **Zero file writes** under targetDir.

## Errors

- Creating Backup paths fails the leaf.
- Non-zero exit fails.
- Missing symlink dry-run messaging fails.

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
	if pathLexists(req.BackupLinkPath) {
		t.Fatalf("dry-run must not create %s", req.BackupLinkPath)
	}
	if pathLexists(req.BackupFilePath) {
		t.Fatalf("dry-run must not create %s", req.BackupFilePath)
	}

	combined := combinedOut(resp)
	if strings.Contains(combined, "would symlink") {
		return
	}
	if strings.Contains(combined, "dry-run") && strings.Contains(combined, "symlink") {
		return
	}
	t.Fatalf("dry-run: prefer log with would symlink\nstdout:\n%s\nstderr:\n%s",
		resp.Stdout, resp.Stderr)
}
```
