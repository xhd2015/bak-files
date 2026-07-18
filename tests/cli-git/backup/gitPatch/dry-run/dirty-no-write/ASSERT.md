# gitPatch --dry-run dirty → no patch, exit 0

## Expected

- Exit code **0**.
- No `*.patch` under **PatchDir** / **TargetDir**.
- TargetDir fingerprint equals **TargetDirBefore**.
- Prefer log contains `dry-run` or `would`.

## Side Effects

- Zero patch writes.
- No remote ops.

## Errors

- Any patch or targetDir churn fails.
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
		t.Fatalf("dry-run must not write patch %s", req.BackupPath)
	}
	if found := findPatchFiles(t, req.TargetDir); len(found) > 0 {
		t.Fatalf("dry-run must not write patches: %v\nstdout:\n%s\nstderr:\n%s",
			found, resp.Stdout, resp.Stderr)
	}
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("dry-run wrote under targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "dry-run") && !strings.Contains(combined, "would") {
		t.Fatalf("dry-run: prefer log with dry-run or would\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
