# Scenario

**Feature**: `.grok/projects` is per-workspace agent cache

```
caller -> Classify(".grok/projects/Users-x/terminals/1.txt")
  -> Rule=.grok/projects, Flags=cache, Owner=grok
```

## Preconditions

- Catalog: `.grok/projects` → Cache, Owner grok.

## Steps

1. Set path under `.grok/projects`.
2. Expect cache flags and owner grok.

## Context

- Workspace-local tool/terminal state; regenerable, not history.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".grok/projects/Users-x/terminals/1.txt"
	return nil
}
```
