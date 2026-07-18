# Scenario

**Feature**: `.wrk` is worktree/workspace cache without owner

```
caller -> Classify(".wrk/something")
  -> Rule=.wrk, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.wrk` → Cache, Owner none.

## Steps

1. Set path under `.wrk`.
2. Expect cache flags and empty owner.

## Context

- Worktree roots are regenerable working/cache state, not user history.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".wrk/something"
	return nil
}
```
