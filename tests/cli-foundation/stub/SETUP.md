# Scenario

**Feature**: known subcommands without help are stubs (not implemented)

```
# operator runs a known command without --help
operator -> bak-files backup|restore|list
  -> not implemented signal (exit non-zero)
```

## Preconditions

- No help flags.
- P1 does not require real backup/restore/list behavior.

## Steps

1. Child sets `Args` to a single command name.
2. Assert non-zero exit and not-implemented wording on stdout or stderr.

## Context

- Prefer non-zero exit so later phases can replace stubs with success paths cleanly.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
