# Scenario

**Feature**: global basename excludes (`*.tmp`) still work with dots discovery

```
# operator
write home/.myapp/keep.txt = "keep-me\n"
write home/.myapp/noise.tmp = "noise\n"
config: files={}, global.excludes=["*.tmp"]
operator -> bak-files backup --config …
  -> keep.txt under targetDir; noise.tmp MUST NOT
  -> exit 0
```

## Preconditions

- `.myapp` is not pathflag DefaultSkip (ordinary owner-less path).
- Basename exclude is step 4 of skip policy after pathflag.

## Steps

1. setupDotsWorld emptyFilesConfigBasenameExcludes, real backup.
2. Plant keep and noise; record paths.

## Context

- Ensures dots feature does not break existing exclude semantics.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfigBasenameExcludes(), "backup", false)
	writeFile(t, filepath.Join(home, ".myapp", "keep.txt"), "keep-me\n")
	writeFile(t, filepath.Join(home, ".myapp", "noise.tmp"), "noise\n")
	req.KeepBackupPath = mappingBackup(req.TargetDir, "alice", ".myapp", "keep.txt")
	req.KeepContent = "keep-me\n"
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".myapp", "noise.tmp")
	return nil
}
```
