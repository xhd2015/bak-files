# Scenario

**Feature**: `.LOG` does not match case-sensitive log suffix

```
caller -> Classify("projects/myapp/debug.LOG")
  -> zero Result (no **/*.log rule)
```

## Preconditions

- Basename ends with `.LOG` (uppercase), not `.log`.

## Steps

1. Set path ending in `.LOG`.
2. Expect no log rule (zero Result unless other rules hit — they should not).

## Context

- Suffix match is case-sensitive per catalog design.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "projects/myapp/debug.LOG"
	return nil
}
```
