# Scenario

**Feature**: invalid JSON config fails non-zero

```
# operator provides a non-JSON config file
operator -> bak-files list --config broken.json
  -> parse/config error (exit non-zero)
```

## Preconditions

- File exists but content is not valid JSON.
- Env has HOME and WORKING_ROLE.

## Steps

1. Write broken JSON to WorkDir.
2. Run `list --config` pointing at that file.

## Context

- TS reference: `parse config: …` on JSON.parse failure.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	cfg := filepath.Join(dir, "broken.json")
	writeFile(t, cfg, "{ this is not valid json near\n")
	req.WorkDir = dir
	req.Env = listEnv(home, "alice")
	req.Args = []string{"list", "--config", cfg}
	return nil
}
```
