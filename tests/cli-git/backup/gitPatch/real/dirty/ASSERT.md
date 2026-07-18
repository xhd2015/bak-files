# Real gitPatch dirty → patch under mapping, exit 0

## Expected

- Exit code **0**.
- Under **PatchDir** (`files/HOME/repo`) there is at least one `*.patch` file,
  **or** **BackupPath** (`…/worktree.patch`) exists.
- Patch content **looksLikeDiff** (`diff --git` and/or unified `---`/`+++`/`@@`).
- Prefer content reflects the dirty change (contains `patched-line` or a `+`
  line) — soft if format uses binary/encoding, but plain text README must
  appear as a text hunk in standard `git diff HEAD`.

## Side Effects

- Creates targetDir / mapping directory as needed.
- Does not require network remotes.
- Does not need to write bak.stats for pure gitPatch MVP (optional).

## Errors

- No patch file fails.
- Non-diff content fails.
- Non-zero exit fails.

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
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	var patchPath string
	var body string
	if pathExists(req.BackupPath) {
		patchPath = req.BackupPath
		body = readFileOrEmpty(req.BackupPath)
	} else {
		found := findPatchFiles(t, req.PatchDir)
		if len(found) == 0 {
			// Also accept a single non-.patch file that is pure diff at mapping path
			// as files/HOME/repo if it is a file.
			if st, err := os.Stat(req.PatchDir); err == nil && !st.IsDir() {
				patchPath = req.PatchDir
				body = readFileOrEmpty(req.PatchDir)
			} else {
				t.Fatalf("no .patch under %s and BackupPath missing %s\ntarget tree:\n%s\nstdout:\n%s\nstderr:\n%s",
					req.PatchDir, req.BackupPath, treeFingerprint(t, req.TargetDir), resp.Stdout, resp.Stderr)
			}
		} else {
			patchPath = found[0]
			body = readFileOrEmpty(patchPath)
		}
	}
	if !looksLikeDiff(body) {
		t.Fatalf("patch %s does not look like a diff:\n%s\nstdout:\n%s\nstderr:\n%s",
			patchPath, body, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(body, "patched-line") && !strings.Contains(body, "+patched") {
		// git diff usually includes the new line with a leading +.
		if !strings.Contains(body, "+") {
			t.Fatalf("patch %s missing dirty content markers:\n%s", patchPath, body)
		}
	}
	_ = filepath.Dir(patchPath)
}
```
