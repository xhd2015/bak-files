# Scenario

**Feature**: backup/restore `--help` without config or env

```
# help skips config and does not write trees
operator -> bak-files backup|restore --help
  -> Usage on stdout (exit 0)
```

## Preconditions

- Minimal env (PATH only).
- No WorkDir fixture required.

## Steps

1. Leaves set Args to `{backup|restore, --help}`.

## Context

- Help must remain green when backup/restore bodies are implemented (not only stubs).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Help needs no HOME/config; force minimal env so leaves cannot leak env.
	req.Env = minimalEnv()
	req.WorkDir = ""
	return nil
}
```
