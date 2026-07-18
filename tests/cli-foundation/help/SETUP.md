# Scenario

**Feature**: multi-level help exits 0 with Usage text

```
# operator asks for help at root or subcommand
operator -> bak-files [cmd] --help|-h | (empty)
  -> Usage on stdout (exit 0)
```

## Preconditions

- `Mode` is `cli`.
- Help paths must not require config or env.

## Steps

1. Grouping marks this branch as help-only (exit 0 expected at leaves).
2. Child sets concrete `Args`.

## Context

- Root help lists `backup`, `restore`, `list`.
- Subcommand help is command-specific and does not need to re-list all peers.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
