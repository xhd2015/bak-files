# Scenario

**Feature**: `restore --help` shows command usage

```
# operator asks for restore help
operator -> bak-files restore --help
  -> restore Usage on stdout (exit 0)
```

## Preconditions

- No config required for help.

## Steps

1. Set `Args` to `{"restore", "--help"}`.

## Context

- Help exits 0 even while the restore body is a stub.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"restore", "--help"}
	return nil
}
```
