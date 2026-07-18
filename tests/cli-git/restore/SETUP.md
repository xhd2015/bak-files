# Scenario

**Feature**: `bak-files restore` after a gitTree-style backup seed

```
# restore subcommand
operator -> bak-files restore --config …
  -> plain copy targetDir → source path
```

## Preconditions

- Subcommand is **restore**.
- Advanced git checkout/FF restore is **out of scope**; plain content restore only.

## Steps

1. Leaf seeds targetDir and missing/overwritten source, then restore.

## Context

- Minimal P5 restore happy path.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"restore"}
	}
	return nil
}
```
