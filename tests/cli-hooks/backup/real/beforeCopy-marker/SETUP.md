# Scenario

**Feature**: real backup runs `beforeCopy` (marker file appears)

```
# operator
write home/notes.txt = "hooked-body\n"
config files["~/notes.txt"] = {
  "file": "~/notes.txt",
  "beforeCopy": "touch '<WorkDir>/markers/beforeCopy.ran'"
}
operator -> bak-files backup --config …
  -> markers/beforeCopy.ran exists
  -> files/HOME/notes.txt == "hooked-body\n"
  -> exit 0
```

## Preconditions

- Source `~/notes.txt` exists; marker path under WorkDir does not exist before run.
- Safe script: `touch` only.

## Steps

1. initHookWorld; create source file; config with beforeCopy + file.
2. Record MarkerPath, BackupPath, Content.
3. Args: real `backup --config …`.

## Context

- P4 exit criterion: real run produces hook side effect (marker).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, home := initHookWorld(t, req)

	const body = "hooked-body\n"
	src := filepath.Join(home, "notes.txt")
	writeFile(t, src, body)

	markerDir := filepath.Join(workDir, "markers")
	if err := os.MkdirAll(markerDir, 0o755); err != nil {
		t.Fatalf("mkdir markers: %v", err)
	}
	marker := filepath.Join(markerDir, "beforeCopy.ran")
	// Ensure absent before run.
	_ = os.Remove(marker)

	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/notes.txt": map[string]any{
			"file":       "~/notes.txt",
			"beforeCopy": "touch " + shellQuote(marker),
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = src
	req.BackupPath = filepath.Join(req.TargetDir, "HOME", "notes.txt")
	req.Content = body
	req.MarkerPath = marker
	req.MarkerBefore = ""
	req.Args = []string{"backup", "--config", cfgPath}
	return nil
}
```
