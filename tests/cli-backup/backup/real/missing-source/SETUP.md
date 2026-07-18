# Scenario

**Feature**: missing backup source is skipped (INFO), exit 0

```
# operator
config lists ~/does-not-exist.txt; file not created
operator -> bak-files backup --config …
  -> skip/info about missing path
  -> exit 0 (prefer TS-compatible; do not fail whole run)
  -> no spurious backup file for that path (or empty store OK)
```

## Preconditions

- Config entry path does not exist under HOME.
- Env valid (HOME, WORKING_ROLE set).

## Steps

1. Write `missingSourceConfig()` only (no source file).
2. Args: `backup --config <path>`.
3. Record expected BackupPath that must not be created as a real content file
   (or may be absent).

## Context

- Prefer messages mentioning skip/missing/not found on stdout or stderr;
  exact wording is soft if exit 0 and no invented content.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir home: %v", err)
	}
	cfgPath := filepath.Join(dir, "bak.config.json")
	writeFile(t, cfgPath, missingSourceConfig())
	target := filepath.Join(dir, "files")
	req.WorkDir = dir
	req.Env = bakEnv(home, "alice")
	req.TargetDir = target
	req.SourcePath = filepath.Join(home, "does-not-exist.txt")
	req.BackupPath = filepath.Join(target, "HOME", "does-not-exist.txt")
	req.Args = []string{"backup", "--config", cfgPath}
	return nil
}
```
