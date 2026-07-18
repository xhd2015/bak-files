# Scenario

**Feature**: `bak-files backup --help` documents dotfile flags

```
# operator
operator -> bak-files backup --help
  -> Usage … --no-dot-files … --include … --exclude … (exit 0)
```

## Preconditions

- Args: `backup --help`.

## Steps

1. Set Args to `backup --help`.
2. Env already minimal from parent.

## Context

- Requirement: update backup help for `--no-dot-files`, `--include`, `--exclude`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"backup", "--help"}
	return nil
}
```
