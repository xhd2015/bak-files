# Scenario

**Feature**: `"~": [".bashrc"]` must not full-home-copy `Library`

```
# operator
HOME has .bashrc and Library/Android/sdk/x
files: { "~": [".bashrc"] }  # dots default on (includeDotFiles omitted)
operator -> bak-files backup --config …
  -> files/HOME/alice/.bashrc == "rc-body\n"
  -> files/HOME/alice/Library/... MUST NOT exist
  -> exit 0
```

## Preconditions

- `tildeArrayConfig(['.bashrc'])`; poison tree `Library/Android/sdk/x` under HOME.
- Mapping still `~` → `HOME/$WORKING_ROLE`.

## Steps

1. setupDotsWorld with tildeArrayConfig, real backup.
2. Plant `.bashrc` and `Library/Android/sdk/x`.
3. KeepBackupPath = mapped `.bashrc`; ExcludedBackupPath = mapped Library file.

## Context

- Primary bug: bare `"~"` becomes Source=$HOME and walks entire home including Library.
- Array value must expand to `$HOME/.bashrc` job (and discovery may add other dots, still never Library).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, tildeArrayConfig(`[".bashrc"]`), "backup", false)
	const body = "rc-body\n"
	writeFile(t, filepath.Join(home, ".bashrc"), body)
	writeFile(t, filepath.Join(home, "Library", "Android", "sdk", "x"), "poison-sdk\n")
	req.SourcePath = filepath.Join(home, ".bashrc")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".bashrc")
	req.KeepContent = body
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", "Library", "Android", "sdk", "x")
	return nil
}
```
