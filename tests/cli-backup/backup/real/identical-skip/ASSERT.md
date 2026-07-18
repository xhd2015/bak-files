# Identical content backup → bytes unchanged, exit 0

## Expected

- Exit code **0**.
- **BackupPath** content still `same-bytes\n`.
- TargetDir fingerprint equals **TargetDirBefore** (no extra files; same bodies).

## Side Effects

- Prefer no-op rewrite; tree fingerprint must match pre-run snapshot.

## Errors

- Content or tree drift fails.
- Non-zero exit fails.

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
	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("backup content changed:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			got, req.Content, resp.Stdout, resp.Stderr)
	}
	after := treeFingerprint(t, req.TargetDir)
	if after != req.TargetDirBefore {
		t.Fatalf("targetDir tree changed after identical backup\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, after, resp.Stdout, resp.Stderr)
	}
}
```
