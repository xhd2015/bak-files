# Scenario

**Feature**: valid config + env → list prints mapping paths

```
# happy path
operator -> bak-files list [--config FILE] [NAMES...]
  -> load + validate + resolve mapping
  -> stdout mappingPath lines (exit 0)
```

## Preconditions

- Leaves set HOME, WORKING_ROLE, fixture files, and Args.
- Fixture mapping (see `standardListFixture`):
  - `~/Scripts` → `HOME/Scripts`
  - `~` → `HOME/$WORKING_ROLE`
- With `WORKING_ROLE=alice`, expected paths for the standard fixture:
  - `HOME/alice/.bashrc`
  - `HOME/Scripts/tool.sh`

## Steps

1. Grouping documents success contract for descendants.

## Context

- Order: preserve `files` key order from the JSON fixture (not random map order).
- Mapping applies longest/declared prefixes so `~/Scripts/...` uses Scripts
  mapping, not only `~`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Success leaves always exercise list (children set config path / names).
	if len(req.Args) == 0 {
		req.Args = []string{"list"}
	}
	return nil
}
```
