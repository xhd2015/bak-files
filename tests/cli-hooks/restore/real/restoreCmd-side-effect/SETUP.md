# Scenario

**Feature**: real restore runs `restoreCmd` (side-effect file appears)

```
# operator
seed files/HOME/notes.txt = "from-backup\n"
dest home/notes.txt absent
config files["~/notes.txt"] = {
  "file": "~/notes.txt",
  "restoreCmd": "touch '<WorkDir>/markers/restoreCmd.ran'"
}
# (implementations may use beforeRestore instead; restoreCmd is preferred here)
operator -> bak-files restore --config …
  -> markers/restoreCmd.ran exists
  -> home/notes.txt == "from-backup\n"  (if restore still copies file)
  -> exit 0
```

## Preconditions

- Backup store pre-seeded; destination absent; marker absent.
- Safe script: `touch` only.

## Steps

1. initHookWorld; seed BackupPath; config with restoreCmd.
2. Record SideEffectPath, SourcePath, Content.
3. Args: real `restore --config …`.

## Context

- P4 optional easy path: restoreCmd side effect on real restore.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, home := initHookWorld(t, req)

	const body = "from-backup\n"
	src := filepath.Join(home, "notes.txt")
	backup := filepath.Join(req.TargetDir, "HOME", "notes.txt")
	writeFile(t, backup, body)
	_ = os.Remove(src)

	markerDir := filepath.Join(workDir, "markers")
	if err := os.MkdirAll(markerDir, 0o755); err != nil {
		t.Fatalf("mkdir markers: %v", err)
	}
	marker := filepath.Join(markerDir, "restoreCmd.ran")
	_ = os.Remove(marker)

	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/notes.txt": map[string]any{
			"file":       "~/notes.txt",
			"restoreCmd": "touch " + shellQuote(marker),
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = src
	req.BackupPath = backup
	req.Content = body
	req.SideEffectPath = marker
	req.SideEffectBefore = ""
	req.Args = []string{"restore", "--config", cfgPath}
	return nil
}
```
