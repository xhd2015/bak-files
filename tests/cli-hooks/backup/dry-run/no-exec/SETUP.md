# Scenario

**Feature**: dry-run does not execute `beforeCopy` or `cmd` (no marker, no file)

```
# operator
write home/notes.txt = "dry-hook-src\n"
config:
  files["~/generated.txt"] = { "cmd": "printf 'should-not-write\\n'" }
  files["~/notes.txt"] = {
    "file": "~/notes.txt",
    "beforeCopy": "touch '<WorkDir>/markers/dry.ran'"
  }
operator -> bak-files backup --config … --dry-run
  -> exit 0
  -> markers/dry.ran absent
  -> files/HOME/generated.txt absent
  -> prefer log contains dry-run or would
```

## Preconditions

- Config includes **both** a `cmd` entry and a `beforeCopy` entry that would
  leave fingerprints on a real run.
- TargetDir absent before run; marker absent.

## Steps

1. initHookWorld; write source for the file entry; config with cmd + beforeCopy.
2. Args include `--dry-run`.
3. Record MarkerPath, BackupPath (cmd target), TargetDirBefore.

## Context

- P4 exit criterion: dry-run does not run script and does not write.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, home := initHookWorld(t, req)

	const body = "dry-hook-src\n"
	src := filepath.Join(home, "notes.txt")
	writeFile(t, src, body)

	markerDir := filepath.Join(workDir, "markers")
	if err := os.MkdirAll(markerDir, 0o755); err != nil {
		t.Fatalf("mkdir markers: %v", err)
	}
	marker := filepath.Join(markerDir, "dry.ran")
	_ = os.Remove(marker)

	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/generated.txt": map[string]any{
			"cmd": "printf 'should-not-write\\n'",
		},
		"~/notes.txt": map[string]any{
			"file":       "~/notes.txt",
			"beforeCopy": "touch " + shellQuote(marker),
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = src
	req.MarkerPath = marker
	req.MarkerBefore = ""
	// Primary cmd-generated path that must stay absent.
	req.BackupPath = filepath.Join(req.TargetDir, "HOME", "generated.txt")
	req.Content = "should-not-write\n"
	req.TargetDirBefore = treeFingerprint(t, req.TargetDir)
	req.Args = []string{"backup", "--config", cfgPath, "--dry-run"}
	return nil
}
```
