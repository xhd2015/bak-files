# Scenario

**Feature**: `.commandcode/projects` is Commandcode per-workspace project cache

```
caller -> Classify(".commandcode/projects/users-x/y")
  -> Rule=.commandcode/projects, Flags=cache, Owner=commandcode
```

## Preconditions

- Catalog: `.commandcode/projects` → Cache, owner commandcode.
- Whole `.commandcode` root is not a catalog rule.

## Steps

1. Set path under `.commandcode/projects`.
2. Expect cache flags and owner commandcode.

## Context

- Fine prefix only; project-local regenerable state.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".commandcode/projects/users-x/y"
	return nil
}
```
