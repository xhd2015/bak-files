# Scenario

**Feature**: `bak-files backup --help` shows backup usage

```
# operator
operator -> bak-files backup --help
  -> Usage … backup … (exit 0)
```

## Preconditions

- Args: `backup --help`.

## Steps

1. Set Args to `backup --help`.
2. Env already minimal from parent.

## Context

- Prefer documenting `--config` and `--dry-run` when backup is real.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"backup", "--help"}
	return nil
}
```
