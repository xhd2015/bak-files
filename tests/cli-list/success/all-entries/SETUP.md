# Scenario

**Feature**: list all resolved mapping paths via `--config`

```
# operator
operator -> bak-files list --config bak.config.json
  env HOME, WORKING_ROLE=alice
  ->
HOME/alice/.bashrc
HOME/Scripts/tool.sh
```

## Preconditions

- Standard fixture; both entries listed; no NAMES filter.

## Steps

1. Write `standardListFixture()` under WorkDir.
2. Env: HOME=<workdir>/home, WORKING_ROLE=alice.
3. Args: `list --config <path>`.

## Context

- Primary exit criterion for P2: fixture + env prints expected names.

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
	req.Args = []string{"list", "--config", cfg}
	return nil
}
```
