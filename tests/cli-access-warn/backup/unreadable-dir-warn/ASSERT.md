# unreadable-dir-warn → warning: skip, exit 0, bashrc copied

## Expected

- Exit code **0**.
- Combined stdout/stderr contains **`warning:`** and **`skip`** and a path
  fragment for **`.private`** (engine format: `warning: skip <path>: …`).
- **BackupPath** (`.bashrc`) exists with **Content**.
- Combined logs must **not** show a fatal whole-run abort token **`Error:`**
  that fails the process (soft: if `error:` appears, exit must still be 0 and
  bashrc present — primary signal is exit 0 + warning skip).
- Prefer **ExcludedBackupPath** (secret) absent under targetDir.

## Side Effects

- May create empty or partial `.private` mapping dir; nested secret not required.
- Restores unreadable dir mode at end for cleanup.

## Errors

- Non-zero exit fails.
- Missing warning/skip/.private tokens fails.
- Missing bashrc fails.

## Exit Code

- **0**

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// Always restore mode so TempDir cleanup succeeds even on assert failure paths.
	defer restoreUnreadable(req.UnreadableDir)

	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	combined := combinedOut(resp)
	if !strings.Contains(combined, "warning:") && !strings.Contains(combined, "warning") {
		t.Fatalf("expected warning: for unreadable dir\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, "skip") {
		t.Fatalf("expected skip in warning path\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, ".private") && !strings.Contains(combined, "private") {
		t.Fatalf("expected .private path in skip warning\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}

	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("bashrc backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("secret under unreadable dir must not be backed up: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
}
```
