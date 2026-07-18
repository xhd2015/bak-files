# Scenario

**Feature**: backup `cmd` entry writes mapping-path file from command stdout

```
# operator
HOME=…/home  WORKING_ROLE=alice
config files["~/generated.txt"] = { "cmd": "printf 'hello-from-cmd\\n'" }
operator -> bak-files backup --config bak.config.json
  -> files/HOME/generated.txt == "hello-from-cmd\n"
  -> exit 0
```

## Preconditions

- Config object entry with **`cmd`** only (no source file required).
- Mapping `~` → `HOME`; expected BackupPath `WorkDir/files/HOME/generated.txt`.
- Env: HOME, WORKING_ROLE=alice.

## Steps

1. initHookWorld; write config with cmd entry.
2. Args: `backup --config <path>` (no dry-run).
3. Record BackupPath, Content for Assert.

## Context

- P4 exit criterion: real run produces file from cmd stdout.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	workDir, _ := initHookWorld(t, req)

	const body = "hello-from-cmd\n"
	cfg := baseConfigSkeleton()
	cfg["files"] = map[string]any{
		"~/generated.txt": map[string]any{
			"cmd": "printf 'hello-from-cmd\\n'",
		},
	}
	cfgPath := writeConfigJSON(t, workDir, cfg)

	req.BackupPath = filepath.Join(req.TargetDir, "HOME", "generated.txt")
	req.Content = body
	req.Args = []string{"backup", "--config", cfgPath}
	return nil
}
```
