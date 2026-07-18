# Scenario

**Feature**: bare `backup` without config fails opening config (backup is implemented)

```
# operator runs bak-files backup with no bak.config.json in cwd
operator -> bak-files backup
  -> stderr Error: open config … bak.config.json (exit non-zero)
```

## Preconditions

- Args: `backup` only.
- Cwd is module root (default Run cwd), which has **no** `bak.config.json` for this leaf.
- P3 implemented real `backup`; missing config is a normal failure path, not a stub.

## Steps

1. Set `Args` to `{"backup"}`.
2. Run without creating a config file.

## Context

- Same missing-config failure shape as bare `list` / `restore`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"backup"}
	return nil
}
```
