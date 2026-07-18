# Scenario

**Feature**: explicit `~/Scripts` still works; non-listed `Library` is not a job

```
# operator
files: { "~/Scripts": true }  # no bare "~" key; dots default on
HOME: Scripts/tool.sh, Library/x, (optional .bashrc for discovery — not asserted)
operator -> bak-files backup --config …
  -> files/HOME/alice/Scripts/tool.sh == "echo ok\n"
  -> Library under store MUST NOT exist
  -> exit 0
```

## Preconditions

- `explicitScriptsConfig`; Scripts is ordinary `~/path` key (ExpandPath unchanged).
- Library is poison non-dot content that must never be pulled without a job.

## Steps

1. setupDotsWorld explicitScriptsConfig, real backup.
2. Plant Scripts and Library; set Keep / Excluded paths.

## Context

- Regression: fixing `"~"` array must not break explicit `~/Scripts` keys.
- Library absence proves whitelist-only (no accidental home-root job).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, explicitScriptsConfig(), "backup", false)
	const body = "echo ok\n"
	writeFile(t, filepath.Join(home, "Scripts", "tool.sh"), body)
	writeFile(t, filepath.Join(home, "Library", "x"), "lib-poison\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", "Scripts", "tool.sh")
	req.KeepContent = body
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", "Library", "x")
	return nil
}
```
