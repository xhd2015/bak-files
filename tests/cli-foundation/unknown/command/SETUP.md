# Scenario

**Feature**: unknown command name is rejected

```
# operator runs bak-files not-a-real-command
operator -> bak-files not-a-real-command
  -> stderr with Error: or bak-files: (exit non-zero)
```

## Preconditions

- First arg is not a known subcommand.

## Steps

1. Set `Args` to `{"not-a-real-command"}`.

## Context

- Distinguishes unknown routing from stub “not implemented” for known commands.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"not-a-real-command"}
	return nil
}
```
