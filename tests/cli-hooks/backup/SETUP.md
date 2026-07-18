# Scenario

**Feature**: `bak-files backup` runs entry `cmd` / `beforeCopy` / `afterCopy`
(or dry-runs without executing them)

```
# backup hooks pipeline
operator -> bak-files backup --config … [--dry-run]
  -> (real) spawn shell for cmd/beforeCopy; write mapping files
  -> (dry-run) no shell for hooks/cmd; no generated/marker files
```

## Preconditions

- Leaves supply WorkDir, config with hook/cmd entries, HOME, WORKING_ROLE,
  and Args starting with `backup`.

## Steps

1. Grouping documents real vs dry-run branches under backup.

## Context

- Mapping: `~` → `HOME`; targetDir `./files`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	}
	return nil
}
```
