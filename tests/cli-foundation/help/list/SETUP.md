# Scenario

**Feature**: `list --help` shows command usage

```
# operator asks for list help
operator -> bak-files list --help
  -> list Usage on stdout (exit 0)
```

## Preconditions

- No config required for help.

## Steps

1. Set `Args` to `{"list", "--help"}`.

## Context

- Help exits 0 even while the list body is a stub.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"list", "--help"}
	return nil
}
```
