# Scenario

**Feature**: missing `--config` path fails non-zero

```
# operator points at a non-existent config
operator -> bak-files list --config /no/such/bak.config.json
  -> error on stderr (exit non-zero)
```

## Preconditions

- HOME and WORKING_ROLE may be set (failure is file-not-found, not env).
- Config path must not exist.

## Steps

1. Create empty WorkDir (no config files).
2. Set Args to `list --config <workdir>/missing-bak.config.json`.
3. Set Env with HOME + WORKING_ROLE so the only failure is missing file.

## Context

- Default-config missing is similar; this leaf fixes `--config` explicitly.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	missing := filepath.Join(dir, "missing-bak.config.json")
	req.WorkDir = dir
	req.Env = listEnv(home, "alice")
	req.Args = []string{"list", "--config", missing}
	return nil
}
```
