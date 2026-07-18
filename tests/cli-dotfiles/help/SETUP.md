# Scenario

**Feature**: help documents dotfile CLI flags without config or env

```
# help skips config and does not write trees
operator -> bak-files backup --help
  -> Usage on stdout (exit 0); documents new flags
```

## Preconditions

- Minimal env (PATH only).
- No WorkDir fixture required.

## Steps

1. Leaves set Args to `backup --help`.

## Context

- Help must not require HOME or bak.config.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = minimalEnv()
	req.WorkDir = ""
	return nil
}
```
