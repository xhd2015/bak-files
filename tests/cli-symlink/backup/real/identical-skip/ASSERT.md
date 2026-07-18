# Second backup → exit 0, symlink + file unchanged

## Expected

- Exit code **0**.
- **BackupFilePath** still equals **Content**.
- **BackupLinkPath** still a symlink with **LinkTarget**.
- Prefer **TargetDir** fingerprint equals **TargetDirBefore** (stable tree).

## Side Effects

- Prefer no-op when link target already matches; must not corrupt link.

## Errors

- Non-zero exit fails.
- Drift of file body or link target fails.
- Tree fingerprint drift fails when TargetDirBefore was set.

## Exit Code

- **0**

```go
import "testing"

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
		t.Fatalf("backup file changed:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			got, req.Content, resp.Stdout, resp.Stderr)
	}
	if !isSymlink(req.BackupLinkPath) {
		t.Fatalf("backup link lost at %s\nstdout:\n%s\nstderr:\n%s",
			req.BackupLinkPath, resp.Stdout, resp.Stderr)
	}
	gotTgt := readlinkOrEmpty(req.BackupLinkPath)
	if gotTgt != req.LinkTarget {
		t.Fatalf("backup symlink target changed:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			gotTgt, req.LinkTarget, resp.Stdout, resp.Stderr)
	}

	if req.TargetDirBefore != "" {
		after := treeFingerprint(t, req.TargetDir)
		if after != req.TargetDirBefore {
			t.Fatalf("targetDir tree changed after identical backup\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
				req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
		}
	}
}
```
