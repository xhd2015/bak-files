# Scenario

**Feature**: bare `restore` without config fails opening config (restore is implemented)

```
# operator runs bak-files restore with no bak.config.json in cwd
operator -> bak-files restore
  -> stderr Error: open config … bak.config.json (exit non-zero)
```

## Preconditions

- Args: `restore` only.
- Cwd is module root (default Run cwd), which has **no** `bak.config.json` for this leaf.
- P3 implemented real `restore`; missing config is a normal failure path, not a stub.

## Steps

1. Set `Args` to `{"restore"}`.
2. Run without creating a config file.

## Context

- Same missing-config failure shape as bare `list` / `backup`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"restore"}
	return nil
}
```
