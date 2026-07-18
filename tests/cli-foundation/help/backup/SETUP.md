# Scenario

**Feature**: `backup --help` shows command usage

```
# operator asks for backup help
operator -> bak-files backup --help
  -> backup Usage on stdout (exit 0)
```

## Preconditions

- No config required for help.

## Steps

1. Set `Args` to `{"backup", "--help"}`.

## Context

- Help for the stub command must exit 0 even though the command body is not implemented.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"backup", "--help"}
	return nil
}
```
