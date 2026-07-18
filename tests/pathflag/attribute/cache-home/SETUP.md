# Scenario

**Feature**: `.cache` is temporary application cache without owner

```
caller -> Classify(".cache/something")
  -> Rule=.cache, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.cache` → Cache, Owner none.

## Steps

1. Set path under `.cache`.
2. Expect cache flags and empty owner.

## Context

- Generic XDG-style cache is not bound to a tool owner.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".cache/something"
	return nil
}
```
