# File symlink backup → symlink recreated, target file content copied

## Expected

- Exit code **0**.
- **BackupFilePath** (`…/target.txt`) equals **Content**.
- **BackupLinkPath** (`…/link.txt`) is a **symlink** with Readlink **LinkTarget**
  (`target.txt`).
- Source link remains a symlink with the same target string.

## Side Effects

- Mapping path under `HOME/alice/Scripts/` receives both entries.

## Errors

- Non-zero exit fails.
- If link path is a regular file with body instead of a symlink, fail
  (must preserve as link).

## Exit Code

- **0**

```go
import (
	"os"
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

	got := readFileOrEmpty(req.BackupFilePath)
	if got != req.Content {
		t.Fatalf("backup target file %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupFilePath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	if !isSymlink(req.BackupLinkPath) {
		// Fail hard if implementation Open-copied through the link as a plain file.
		body := readFileOrEmpty(req.BackupLinkPath)
		t.Fatalf("backup path %s: want symlink (target %q); got non-link (body %q)\nstdout:\n%s\nstderr:\n%s",
			req.BackupLinkPath, req.LinkTarget, body, resp.Stdout, resp.Stderr)
	}
	gotTgt := readlinkOrEmpty(req.BackupLinkPath)
	if gotTgt != req.LinkTarget {
		t.Fatalf("backup symlink target:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			gotTgt, req.LinkTarget, resp.Stdout, resp.Stderr)
	}

	if !isSymlink(req.SourceLinkPath) {
		t.Fatalf("source link missing/altered: %s", req.SourceLinkPath)
	}
	if readlinkOrEmpty(req.SourceLinkPath) != req.LinkTarget {
		t.Fatalf("source link target changed")
	}
	src, err := os.ReadFile(req.SourceFilePath)
	if err != nil {
		t.Fatalf("source target missing: %v", err)
	}
	if string(src) != req.Content {
		t.Fatalf("source target altered: got %q want %q", src, req.Content)
	}
}
```
