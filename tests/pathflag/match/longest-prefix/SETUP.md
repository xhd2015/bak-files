# Scenario

**Feature**: longest catalog prefix wins for OpenCode log paths

```
# shorter owner/path prefixes exist, but log rule is longer
caller -> Classify(".local/share/opencode/log/app.log")
  -> Rule=.local/share/opencode/log, Flags=logs, Owner=opencode
```

## Preconditions

- Catalog includes `.local/share/opencode/log` → Logs.
- Owner prefix `.local/share/opencode` alone is not enough for flags.

## Steps

1. Set path under the log rule.
2. Expect the long log rule, not owner-only.

## Context

- Contrasts with `owner-only/opencode-local-share` (config under same share root).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/share/opencode/log/app.log"
	return nil
}
```
