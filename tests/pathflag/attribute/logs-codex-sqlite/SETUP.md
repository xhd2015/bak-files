# Scenario

**Feature**: `.codex/logs_2.sqlite` is Codex logs DB

```
caller -> Classify(".codex/logs_2.sqlite")
  -> Rule=.codex/logs_2.sqlite, Flags=logs, Owner=codex
```

## Preconditions

- Catalog exact path-prefix rule for `.codex/logs_2.sqlite` → Logs, owner codex.
- Not a glob of all `logs_*.sqlite` (non-goal).

## Steps

1. Set path to the catalogued sqlite log file (exact rule match).
2. Expect logs flag and codex owner.

## Context

- Specific filename only; other codex paths stay owner-only or other rules.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".codex/logs_2.sqlite"
	return nil
}
```
