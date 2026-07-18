# Scenario

**Feature**: failing `beforeCopy` causes non-zero process exit

```
# operator
write home/notes.txt = "should-not-matter\n"
config files["~/notes.txt"] = {
  "file": "~/notes.txt",
  "beforeCopy": "exit 1"
}
operator -> bak-files backup --config …
  -> exit code ≠ 0
  -> (prefer) no successful completion of the run
```

## Preconditions

- Source exists so failure is from the hook, not missing path.
- Safe script: `exit 1` only.

## Steps

1. initHookWorld; source + config with beforeCopy=`exit 1`.
2. Args: real backup.

## Context

- P4 exit criterion: hook failure fails non-zero.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, home := initHookWorld(t, req)

	const body = "should-not-matter\n"
	src := filepath.Join(home, "notes.txt")
	writeFile(t, src, body)

	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/notes.txt": map[string]any{
			"file":       "~/notes.txt",
			"beforeCopy": "exit 1",
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.SourcePath = src
	req.BackupPath = filepath.Join(req.TargetDir, "HOME", "notes.txt")
	req.Content = body
	req.Args = []string{"backup", "--config", cfgPath}
	return nil
}
```
