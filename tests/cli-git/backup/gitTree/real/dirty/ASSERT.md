# Real gitTree dirty → copy under targetDir + bak.stats, exit 0

## Expected

- Exit code **0**.
- **BackupPath** (`files/HOME/repo/README.md`) exists and equals **Content**
  (`dirty-readme\n`) — worktree content, not only committed blob.
- **StatsPath** (`WorkDir/bak.stats`) is valid JSON with key **MappingKey**
  (`HOME/repo`) recording **`hasChange`: true** and a **commitHash** matching
  fixture HEAD (full hash or shared prefix ≥ 7 chars).
- Source repo still contains the dirty README (backup is copy, not clean).

## Side Effects

- Creates `targetDir` tree as needed.
- Writes/updates `bak.stats` in process cwd (WorkDir).
- Must not require network remotes.

## Errors

- Non-zero exit fails.
- Missing BackupPath or wrong content fails.
- Missing/invalid bak.stats or hasChange≠true fails.

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
	got := readFileOrEmpty(req.BackupPath)
	if got != req.Content {
		t.Fatalf("backup file %s:\n got %q\nwant %q\nstdout:\n%s\nstderr:\n%s",
			req.BackupPath, got, req.Content, resp.Stdout, resp.Stderr)
	}
	src, err := os.ReadFile(filepath.Join(req.SourcePath, "README.md"))
	if err != nil {
		t.Fatalf("source README missing after backup: %v", err)
	}
	if string(src) != req.Content {
		t.Fatalf("source altered: got %q want %q", src, req.Content)
	}
	stats := readFileOrEmpty(req.StatsPath)
	if err := statsHasChangeTrue(stats, req.MappingKey, req.CommitHash); err != nil {
		t.Fatalf("bak.stats: %v\ncontent:\n%s\nstdout:\n%s\nstderr:\n%s",
			err, stats, resp.Stdout, resp.Stderr)
	}
	// Prefer some log mention of copy or backup (soft: not required if files prove success).
	_ = strings.TrimSpace(resp.Stdout + resp.Stderr)
}
```
