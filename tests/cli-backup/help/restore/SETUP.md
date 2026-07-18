# Scenario

**Feature**: `bak-files restore --help` shows restore usage

```
# operator
operator -> bak-files restore --help
  -> Usage … restore … (exit 0)
```

## Preconditions

- Args: `restore --help`.

## Steps

1. Set Args to `restore --help`.

## Context

- Prefer documenting `--config` and `--dry-run` when restore is real.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"restore", "--help"}
	return nil
}
```
