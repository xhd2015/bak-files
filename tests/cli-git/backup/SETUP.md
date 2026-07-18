# Scenario

**Feature**: `bak-files backup` with git-backed files entries

```
# backup subcommand
operator -> bak-files backup --config bak.config.json [--dry-run]
  -> gitTree | gitPatch entry paths under local ~/repo
```

## Preconditions

- Subcommand is **backup** (not restore).
- Config and local git fixture prepared by descendant leaves.

## Steps

1. Narrow mode to backup; leaves set Args and fixtures.

## Context

- Top split under root for write direction: filesystem → targetDir / stats / patch.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Grouping: backup command; leaves append flags/config.
	if len(req.Args) == 0 {
		req.Args = []string{"backup"}
	}
	return nil
}
```
