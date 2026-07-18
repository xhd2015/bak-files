# Scenario

**Feature**: root `-h` shows usage

```
# operator runs bak-files -h
operator -> bak-files -h
  -> root Usage + backup/restore/list (exit 0)
```

## Preconditions

- Short help flag only.

## Steps

1. Set `Args` to `{"-h"}`.

## Context

- Same root usage contract as empty argv and `--help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"-h"}
	return nil
}
```
