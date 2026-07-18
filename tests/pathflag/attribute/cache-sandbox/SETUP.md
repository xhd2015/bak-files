# Scenario

**Feature**: `.sandbox` is sandbox working cache without owner

```
caller -> Classify(".sandbox/something")
  -> Rule=.sandbox, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.sandbox` → Cache, Owner none.

## Steps

1. Set path under `.sandbox`.
2. Expect cache flags and empty owner.

## Context

- Sandbox trees are regenerable working/cache state, not user history.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".sandbox/something"
	return nil
}
```
