# Scenario

**Feature**: bare `list` without config fails opening config (list is implemented)

```
# operator runs bak-files list with no bak.config.json in cwd
operator -> bak-files list
  -> stderr Error: open config … bak.config.json (exit non-zero)
```

## Preconditions

- Args: `list` only.
- Cwd is module root (default Run cwd), which has **no** `bak.config.json` for this leaf.
- P2 implemented real `list`; missing config is a normal failure path, not a stub.

## Steps

1. Set `Args` to `{"list"}`.
2. Run without creating a config file.

## Context

- Distinct from `backup`/`restore` stubs: list attempts to open config and reports an open/config error.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"list"}
	return nil
}
```
