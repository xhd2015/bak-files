# Scenario

**Feature**: `.bun` is Bun install cache

```
caller -> Classify(".bun") | Classify(".bun/install/cache")
  -> Rule=.bun, Flags=cache, Owner=bun
```

## Preconditions

- Catalog: `.bun` → Cache, owner bun, reason Bun install cache.

## Steps

1. Set path nested under `.bun`.
2. Expect cache + bun owner.

## Context

- Nested path must still take the `.bun` prefix rule.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".bun/install/cache"
	return nil
}
```
