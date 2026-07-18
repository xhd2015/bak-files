# Scenario

**Feature**: `bak-files restore` may execute `restoreCmd` / `beforeRestore`

```
# restore command hooks
operator -> bak-files restore --config …
  -> (real) restoreCmd / beforeRestore execute → side effects
```

## Preconditions

- Leaves supply WorkDir, seeded backup store when needed, Env, Args with
  `restore`.

## Steps

1. Grouping documents restore real branch for restoreCmd.

## Context

- P4 optional easy path: restoreCmd side effect.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"restore"}
	}
	return nil
}
```
