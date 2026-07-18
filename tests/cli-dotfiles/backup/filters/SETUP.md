# Scenario

**Feature**: CLI `--include` / `--exclude` force keep or skip home-relative paths

```
# force include/exclude override default pathflag / default-include secrets
operator -> bak-files backup … --include PATH --exclude PATH
  -> include keeps; exclude skips; exclude wins if both
```

## Preconditions

- Dots on; leaves pass repeatable `--include` / `--exclude` home-relative paths.

## Steps

1. Leaves build home trees and append filter flags via setupDotsWorld extraArgs.

## Context

- Paths are home-relative (e.g. `.cache`, `.ssh`), not absolute.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Filters branch: leaves set full Args including --include/--exclude.
	_ = req.Args
	return nil
}
```
