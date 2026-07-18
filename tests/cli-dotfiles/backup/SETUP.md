# Scenario

**Feature**: `bak-files backup` discovers home dots and applies walk skip policy

```
# backup pipeline with dots
operator -> bak-files backup --config … [--dry-run] [dot flags…]
  -> targetDir receives keep paths (real) | no writes (dry-run)
  -> pathflag / exclude skips logged or omitted from tree
```

## Preconditions

- Leaves supply WorkDir, config, HOME, WORKING_ROLE, and Args starting with
  `backup`.

## Steps

1. Grouping documents discovery / filters / walk-skip branches.

## Context

- Mapping fixtures: `~` → `HOME/$WORKING_ROLE`; role `alice`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	} else if req.Args[0] != "backup" {
		req.Args = append([]string{"backup"}, req.Args...)
	}
	return nil
}
```
