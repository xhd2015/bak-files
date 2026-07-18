# Scenario

**Feature**: `config.validate[].env` rejects missing extra env

```
# fixture requires W0; HOME and WORKING_ROLE are set; W0 is not
operator -> bak-files list --config …
  -> missing env W0 (exit non-zero)
```

## Preconditions

- Config JSON includes `"validate":[{"env":["HOME","WORKING_ROLE","W0"]}]`.
- Env has HOME and WORKING_ROLE only.

## Steps

1. Write fixture with extra validate env `W0`.
2. Run list with builtin envs only.

## Context

- Distinct from builtin-only check: validates the `validate` array path.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	cfg := filepath.Join(dir, "bak.config.json")
	writeFile(t, cfg, `{
  "validate": [
    {
      "env": ["HOME", "WORKING_ROLE", "W0"]
    }
  ],
  "files": {
    "~/.bashrc": true
  },
  "targetDir": "./files",
  "mapping": {
    "~": "HOME/$WORKING_ROLE"
  }
}
`)
	req.WorkDir = dir
	req.Env = listEnv(home, "alice") // no W0
	req.Args = []string{"list", "--config", cfg}
	return nil
}
```
