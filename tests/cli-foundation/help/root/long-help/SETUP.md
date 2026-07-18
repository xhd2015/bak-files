# Scenario

**Feature**: root `--help` shows usage

```
# operator runs bak-files --help
operator -> bak-files --help
  -> root Usage + backup/restore/list (exit 0)
```

## Preconditions

- Long help flag only; no subcommand.

## Steps

1. Set `Args` to `{"--help"}`.

## Context

- Equivalent listing of commands as empty argv / `-h`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"--help"}
	return nil
}
```
