# Scenario

**Feature**: positional NAMES filter mapping paths

```
# operator passes a NAME matching one mappingPath
operator -> bak-files list --config … HOME/Scripts/tool.sh
  -> only HOME/Scripts/tool.sh
```

## Preconditions

- Same standard fixture and env as all-entries.
- Single exact NAME: `HOME/Scripts/tool.sh`.

## Steps

1. Write standard fixture.
2. Args: `list --config <path> HOME/Scripts/tool.sh`.

## Context

- TS filters on `mappingPath` (not raw config key). Exact match is enough for
  P2; wildcard `*` forms can be added later.

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
	req.Args = []string{"list", "--config", cfg, "HOME/Scripts/tool.sh"}
	return nil
}
```
