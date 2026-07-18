# Scenario

**Feature**: empty argv shows root usage

```
# operator runs bak-files with no arguments
operator -> bak-files
  -> Usage + backup/restore/list on stdout (exit 0)
```

## Preconditions

- `Args` is empty.

## Steps

1. Set `Args` to an empty slice.
2. Run binary with no arguments.

## Context

- Same contract as `--help` / `-h` for listing subcommands.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{}
	return nil
}
```
