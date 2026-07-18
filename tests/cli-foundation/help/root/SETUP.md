# Scenario

**Feature**: root-level help (empty argv, --help, -h)

```
# operator invokes bak-files without a subcommand help path
operator -> bak-files | bak-files --help | bak-files -h
  -> root Usage listing backup, restore, list (exit 0)
```

## Preconditions

- No subcommand token before help flags (except empty argv).

## Steps

1. Child leaf sets `Args` to empty, `{"--help"}`, or `{"-h"}`.

## Context

- Empty argv is treated as help/usage, not as an error.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Root help: no subcommand token; leaf sets empty / --help / -h.
	req.Mode = "cli"
	return nil
}
```
