# Real gitTree clean → skip backup, exit 0

## Expected

- Exit code **0**.
- **BackupPath** does **not** exist (no copy of clean tree), **or**
  targetDir fingerprint equals **TargetDirBefore** (no new files under target).
- **StatsPath**: either missing/empty, **or** if present must **not** record
  `hasChange: true` for **MappingKey** as a result of this skip-only run
  (MVP: prefer file absent — TargetDir was empty and stats empty before).
- Combined stdout/stderr preferably contains skip/clean/INFO wording
  (case-insensitive: `skip` or `clean`).

## Side Effects

- No file churn under targetDir for this entry.
- Does not require remotes.

## Errors

- Non-zero exit fails.
- Creating BackupPath with content fails (would mean backup ran).
- Writing bak.stats with hasChange true for MappingKey fails.

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
		t.Fatalf("exit code: got %d, want 0\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	afterTree := treeFingerprint(t, req.TargetDir)
	if afterTree != req.TargetDirBefore {
		// Stronger: if BackupPath appeared, always fail.
		if pathExists(req.BackupPath) {
			t.Fatalf("clean gitTree must not copy; BackupPath exists %s content=%q\nstdout:\n%s\nstderr:\n%s",
				req.BackupPath, readFileOrEmpty(req.BackupPath), resp.Stdout, resp.Stderr)
		}
		t.Fatalf("clean gitTree must not churn targetDir\nbefore:\n%s\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.TargetDirBefore, afterTree, resp.Stdout, resp.Stderr)
	}
	stats := readFileOrEmpty(req.StatsPath)
	if stats != req.StatsBefore {
		// If stats appeared with hasChange true for our key, fail hard.
		if err := statsHasChangeTrue(stats, req.MappingKey, req.CommitHash); err == nil {
			t.Fatalf("clean skip must not write hasChange:true for %s\nstats:\n%s\nstdout:\n%s\nstderr:\n%s",
				req.MappingKey, stats, resp.Stdout, resp.Stderr)
		}
		// Other stats formats/churn: still fail — clean should leave stats unchanged in MVP.
		t.Fatalf("clean gitTree must leave bak.stats unchanged\nbefore %q\nafter:\n%s\nstdout:\n%s\nstderr:\n%s",
			req.StatsBefore, stats, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "skip") && !strings.Contains(combined, "clean") {
		t.Fatalf("clean gitTree: prefer INFO/skip/clean log\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
