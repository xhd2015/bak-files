# Scenario

**Feature**: `bak-files list --help` shows list usage

```
# operator
operator -> bak-files list --help
  -> Usage … list … (exit 0)
```

## Preconditions

- Args: `list --help`.
- WorkDir may be empty (module root); env need not include HOME.

## Steps

1. Set `Args` to `{"list", "--help"}`.
2. Use minimal env (PATH only) so help cannot accidentally depend on HOME.

## Context

- Prefer help text that documents `--config` and optional NAMES when list is real.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"list", "--help"}
	req.Env = minimalEnv()
	return nil
}
```
