# Scenario

**Feature**: default config path `bak.config.json` in cwd

```
# operator runs list without --config; cwd has bak.config.json
operator -> bak-files list
  cwd = WorkDir with bak.config.json
  -> same mapping paths as --config case
```

## Preconditions

- Config file name exactly `bak.config.json` under WorkDir.
- Args are only `list` (no `--config`).

## Steps

1. Write standard fixture to `WorkDir/bak.config.json`.
2. Env with HOME + WORKING_ROLE=alice.
3. Args: `{"list"}`.

## Context

- Proves default config discovery independent of `--config` flag parsing.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	cfg := filepath.Join(dir, "bak.config.json")
	writeFile(t, cfg, standardListFixture())
	req.WorkDir = dir
	req.Env = listEnv(home, "alice")
	req.Args = []string{"list"}
	return nil
}
```
