# Missing source → skip, exit 0

## Expected

- Exit code **0** (prefer INFO skip missing, not hard failure of the run).
- Combined stdout+stderr preferably mentions skip / missing / not found /
  does-not-exist (soft if exit 0 and no bogus backup content).
- **BackupPath** must not exist as a regular file with invented content
  (absent is OK; empty file not required).

## Side Effects

- No requirement to create targetDir entries for the missing path.

## Errors

- Exit non-zero fails this preferred design (unless implementer documents
  hard-fail — requirement prefers skip + 0).

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
		t.Fatalf("missing source: want exit 0 (skip), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if pathExists(req.BackupPath) {
		// Soft: allow only if somehow empty — still prefer absent
		got := readFileOrEmpty(req.BackupPath)
		if got != "" {
			t.Fatalf("missing source should not produce backup content at %s: %q",
				req.BackupPath, got)
		}
	}
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	// Soft signal: if any output, prefer skip/missing wording; silence is OK
	// as long as exit 0 and no content (some CLIs log only at -v).
	_ = combined
}
```
