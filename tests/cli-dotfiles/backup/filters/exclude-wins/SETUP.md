# Scenario

**Feature**: exclude wins when both `--include` and `--exclude` match

```
# operator
write home/.cache/x = "x\n"
operator -> bak-files backup --config … --include .cache --exclude .cache
  -> .cache/x MUST NOT be under targetDir
  -> exit 0
```

## Preconditions

- Same home-relative path on both flags.

## Steps

1. setupDotsWorld with both flags.
2. Write `.cache/x`; set ExcludedBackupPath.

## Context

- Skip policy: force-exclude after force-include; exclude wins if both.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, home, _ := setupDotsWorld(t, req, emptyFilesConfig(), "backup", false,
		"--include", ".cache", "--exclude", ".cache")
	writeFile(t, filepath.Join(home, ".cache", "x"), "x\n")
	req.ExcludedBackupPath = mappingBackup(req.TargetDir, "alice", ".cache", "x")
	return nil
}
```
