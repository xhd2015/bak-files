# Scenario

**Feature**: list fails clearly when config or env is invalid

```
# failure paths before printing mapping paths
operator -> bak-files list --config …
  -> missing file | bad JSON | missing env
  -> stderr error, exit non-zero, no successful path list
```

## Preconditions

- Leaves provide a WorkDir and controlled Env as needed.

## Steps

1. Grouping documents error-class leaves under `errors/`.

## Context

- Errors must be actionable (mention missing env name, parse/config, or file).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Error leaves always exercise the list subcommand (not help).
	if len(req.Args) == 0 {
		req.Args = []string{"list"}
	}
	return nil
}
```
