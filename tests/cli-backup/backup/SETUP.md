# Scenario

**Feature**: `bak-files backup` copies sources into `targetDir` (or dry-runs)

```
# backup pipeline
operator -> bak-files backup --config … [--dry-run]
  -> targetDir receives files (real) | no writes (dry-run)
```

## Preconditions

- Leaves supply WorkDir, config, HOME, WORKING_ROLE, and Args starting with
  `backup`.

## Steps

1. Grouping documents real vs dry-run branches.

## Context

- Mapping in fixtures: `~` → `HOME`; targetDir `./files`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Ensure subcommand is backup if leaf forgot (leaves should set full Args).
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	} else if req.Args[0] != "backup" && req.Args[0] != "restore" {
		// Prepend backup for partial Args built by intermediate nodes.
		req.Args = append([]string{"backup"}, req.Args...)
	}
	return nil
}
```

