# Scenario

**Feature**: basename ending in `.log` matches log suffix rule

```
caller -> Classify("projects/myapp/debug.log")
  -> Rule=**/*.log, Flags=logs, Owner empty
```

## Preconditions

- No path or segment rule applies; basename ends with `.log` (case-sensitive).

## Steps

1. Set an ordinary project path ending in `.log`.
2. Expect `**/*.log` + Logs.

## Context

- Contrast with path-catalog log dirs (e.g. `.grok/logs`) which set a different Rule.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "projects/myapp/debug.log"
	return nil
}
```
