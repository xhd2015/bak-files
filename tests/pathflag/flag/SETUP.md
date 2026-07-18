# Scenario

**Feature**: Flag.String and Flag.Names formatting helpers

```
# ascending bit-order names: tmp, cache, logs, binary, trash, meta, vendor, history
pathflag.Flag bits -> String() "tmp|cache|..." | Names() []string
# DefaultSkipMask vs BackupSkipMask membership
```

## Preconditions

- Leaves set `Op` to `flag_string` or `flag_names` and `FlagBits` from constants.
- Mask leaf may also assert `BackupSkipMask` / `DefaultSkipMask` constants directly.

## Steps

1. Leaf selects helper op and bits via `pathflag.Flag*` constants.
2. Run calls real `Flag.String` or `Flag.Names`.
3. Assert exact formatting or mask membership.

## Context

- These leaves do not call Classify; they pin Flag helper contract used by Response.Flags.
- `history` is last in bit order; reclaim DefaultSkipMask excludes it; bak BackupSkipMask includes it.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Leaf overrides Op to flag_string or flag_names.
	if req.Op == "" || req.Op == OpClassify {
		req.Op = OpFlagString
	}
	return nil
}
```
