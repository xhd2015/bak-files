# Scenario

**Feature**: `bak-files backup` preserves symlinks under `~/Scripts`

```
# backup pipeline
operator -> bak-files backup --config … [--dry-run]
  -> real: copySymlink / copy files under targetDir
  -> dry-run: would symlink…; no writes
```

## Preconditions

- Leaves supply WorkDir, config, HOME, WORKING_ROLE, and Args starting with
  `backup`.

## Steps

1. Grouping documents real vs dry-run branches.

## Context

- Mapping fixtures: `~` → `HOME/$WORKING_ROLE`; targetDir `./files`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	} else if req.Args[0] != "backup" && req.Args[0] != "restore" {
		req.Args = append([]string{"backup"}, req.Args...)
	}
	return nil
}
```
