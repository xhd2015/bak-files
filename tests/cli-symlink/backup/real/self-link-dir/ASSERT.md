# Self-link-dir backup → exit 0, link preserved, no "is a directory"

## Expected

- Exit code **0**.
- **BackupFilePath** exists as a regular file with **Content**.
- **BackupLinkPath** is a **symlink** whose `Readlink` equals **LinkTarget**.
- Combined stdout/stderr must **not** contain `is a directory` (case-insensitive).
- Source file and source link remain intact (copy, not move).

## Side Effects

- Creates `targetDir` tree as needed under `HOME/alice/Scripts/proj/`.

## Errors

- Non-zero exit fails.
- Missing/wrong file body or non-symlink backup link fails.
- Any `is a directory` log token fails.

## Exit Code

- **0**

```go
import (
	"os"
	"strings"
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
	combined := combinedOut(resp)
	if strings.Contains(combined, "is a directory") {
		t.Fatalf("backup must not report is a directory for self-link-dir\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}

	got := readFileOrEmpty(req.BackupFilePath)
	if got != req.Content {
		t.Fatalf("backup file %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupFilePath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	if !isSymlink(req.BackupLinkPath) {
		t.Fatalf("backup path %s: want symlink\nstdout:\n%s\nstderr:\n%s",
			req.BackupLinkPath, resp.Stdout, resp.Stderr)
	}
	gotTgt := readlinkOrEmpty(req.BackupLinkPath)
	if gotTgt != req.LinkTarget {
		t.Fatalf("backup symlink target:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			gotTgt, req.LinkTarget, resp.Stdout, resp.Stderr)
	}

	src, err := os.ReadFile(req.SourceFilePath)
	if err != nil {
		t.Fatalf("source file missing after backup: %v", err)
	}
	if string(src) != req.Content {
		t.Fatalf("source file altered: got %q want %q", src, req.Content)
	}
	if !isSymlink(req.SourceLinkPath) {
		t.Fatalf("source link missing/altered: %s", req.SourceLinkPath)
	}
}
```
