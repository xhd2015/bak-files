# Scenario

**Feature**: `"$W0/proj": [".vscode"]` backs up only `.vscode`, never `vendor`

```
# operator (production-shaped bug)
W0=…/w0
tree:
  w0/proj/.vscode/settings.json = "vscode-settings\n"
  w0/proj/vendor/poison.go = "vendor-poison\n"
files: { "$W0/proj": [".vscode"] }
mapping: { "$W0": "W" }
includeDotFiles: false
operator -> bak-files backup --config …
  -> files/W/proj/.vscode/settings.json == "vscode-settings\n"
  -> files/W/proj/vendor/poison.go MUST NOT exist
  -> no vendor tree under store
  -> exit 0
```

## Preconditions

- `projectArrayConfig(['.vscode'])`; poison `vendor/poison.go` under `$W0/proj`.
- Mapping `$W0` → `W` so keep path is `files/W/proj/.vscode/settings.json`.

## Steps

1. setupProjectArrayWorld backup real + projectArrayConfig.
2. Plant `.vscode/settings.json` and `vendor/poison.go`.
3. KeepBackupPath / ExcludedBackupPath under TargetDir.

## Context

- Bug today: array ignored for every key except bare `"~"` → full `$W0/proj`
  job → vendor copied into store (vendor poison).
- Fixed: one child job `$W0/proj/.vscode` only.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, w0, _ := setupProjectArrayWorld(t, req, projectArrayConfig(`[".vscode"]`), "backup", false)
	const body = "vscode-settings\n"
	writeFile(t, filepath.Join(w0, "proj", ".vscode", "settings.json"), body)
	writeFile(t, filepath.Join(w0, "proj", "vendor", "poison.go"), "vendor-poison\n")
	req.KeepBackupPath = filepath.Join(req.TargetDir, "W", "proj", ".vscode", "settings.json")
	req.KeepContent = body
	req.ExcludedBackupPath = filepath.Join(req.TargetDir, "W", "proj", "vendor", "poison.go")
	return nil
}
```
