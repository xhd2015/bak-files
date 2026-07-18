# Scenario

**Feature**: missing builtin env `WORKING_ROLE` fails validation

```
# config is valid JSON; WORKING_ROLE unset
operator -> bak-files list --config ok.json
  -> missing env error (exit non-zero)
```

## Preconditions

- Valid standard fixture on disk.
- Env sets HOME but **not** WORKING_ROLE.

## Steps

1. Write `standardListFixture()` to WorkDir.
2. Env: HOME only (+ PATH).
3. Args: `list --config <file>`.

## Context

- TS always `validateEnvs(["WORKING_ROLE", "HOME"])` before config.validate.
- Assert must mention the missing env name clearly.

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
	// HOME set; WORKING_ROLE intentionally omitted.
	req.Env = minimalEnv("HOME=" + home)
	req.Args = []string{"list", "--config", cfg}
	return nil
}
```
