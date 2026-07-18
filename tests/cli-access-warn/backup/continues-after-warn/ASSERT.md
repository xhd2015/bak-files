# continues-after-warn → exit 0; good files under mapping

## Expected

- Exit code **0**.
- **BackupPath** (`.bashrc`) equals **Content**.
- **OtherBackupPath** (`.other`) equals **OtherContent**.
- Combined logs contain **warning** + **skip** (access path still exercised).
- Secret under unreadable dir not required under targetDir.

## Side Effects

- Mapping receives both good files despite one denied sibling job.

## Errors

- Non-zero exit fails.
- Either good file missing/wrong body fails.

## Exit Code

- **0**

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	defer restoreUnreadable(req.UnreadableDir)

	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("bashrc backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	other := readFileOrEmpty(req.OtherBackupPath)
	if other != req.OtherContent {
		t.Fatalf("other backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.OtherBackupPath, other, req.OtherContent, resp.Stdout, resp.Stderr)
	}

	combined := combinedOut(resp)
	if !strings.Contains(combined, "warning") {
		t.Fatalf("expected warning for unreadable sibling\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, "skip") {
		t.Fatalf("expected skip in warn path\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if pathExists(req.ExcludedBackupPath) {
		t.Fatalf("secret must not be backed up: %s", req.ExcludedBackupPath)
	}
}
```
