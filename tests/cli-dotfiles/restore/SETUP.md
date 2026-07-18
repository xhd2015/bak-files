# Scenario

**Feature**: restore applies the same home pathflag skip policy as backup

```
# restore pipeline
operator -> bak-files restore --config … [--dry-run]
  -> walk restore candidates under mapping; skip pathflag DefaultSkipMask
  -> dry-run: log would/skip; no dest writes
```

## Preconditions

- Leaves seed targetDir backup trees and optional dests under simulated HOME.

## Steps

1. Leaves set Args starting with `restore`.

## Context

- Requirement: dry-run restore respects same skip when applicable.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"restore"}
	} else if req.Args[0] != "restore" {
		req.Args = append([]string{"restore"}, req.Args...)
	}
	return nil
}
```
