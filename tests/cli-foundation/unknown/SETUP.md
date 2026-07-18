# Scenario

**Feature**: unknown subcommand fails with a clear error

```
# operator types a garbage command name
operator -> bak-files <unknown>
  -> stderr Error: or bak-files: (exit non-zero)
```

## Preconditions

- Command name is not `backup`, `restore`, or `list`.

## Steps

1. Child sets a deliberately unknown first token.

## Context

- Errors go to stderr with a conventional prefix.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
