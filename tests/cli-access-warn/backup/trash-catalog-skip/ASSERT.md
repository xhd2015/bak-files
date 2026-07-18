# trash-catalog-skip → INFO skip .Trash, bashrc copied, exit 0

## Expected

- Exit code **0**.
- Combined stdout/stderr contains **`INFO: skip`** (or case-insensitive `info:` +
  `skip`) mentioning **`.Trash`** and/or **trash** (reason **macOS trash**).
- **BackupPath** (`.bashrc`) exists with **Content**.
- **ExcludedBackupPath** (mapping `.Trash`) does **not** exist; no nested trash
  content under targetDir either.
- No fatal backup abort solely due to `.Trash`.

## Side Effects

- Creates `targetDir/HOME/alice/.bashrc` as needed.
- Does not copy `.Trash` into targetDir.

## Errors

- Non-zero exit fails.
- Missing INFO skip / trash signal fails.
- Missing bashrc or present `.Trash` under target fails.

## Exit Code

- **0**

```go
import (
	"os"
	"path/filepath"
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

	combined := combinedOut(resp)
	if !strings.Contains(combined, "skip") {
		t.Fatalf("expected skip log for .Trash\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	// Prefer INFO: skip … (macOS trash); accept path + trash wording.
	if !strings.Contains(combined, ".trash") && !strings.Contains(combined, "trash") {
		t.Fatalf("expected .Trash / trash in skip context\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	// Soft prefer INFO line shape used by logSkip.
	raw := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(raw, "INFO: skip") && !strings.Contains(combined, "info: skip") {
		t.Fatalf("prefer INFO: skip for catalog trash\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}

	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("bashrc backup %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	if pathLexists(req.ExcludedBackupPath) {
		t.Fatalf(".Trash must not appear under targetDir: %s\nstdout:\n%s\nstderr:\n%s",
			req.ExcludedBackupPath, resp.Stdout, resp.Stderr)
	}
	// Belt-and-suspenders: no path component .Trash under targetDir.
	if req.TargetDir != "" && pathExists(req.TargetDir) {
		err := filepath.Walk(req.TargetDir, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(p) == ".Trash" {
				t.Fatalf("found .Trash under targetDir: %s", p)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk targetDir: %v", err)
		}
	}
}
```
