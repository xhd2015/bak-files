# Scenario

**Feature**: `.tmp` is generic temporary cache without owner

```
caller -> Classify(".tmp/foo.png")
  -> Rule=.tmp, Flags=tmp|cache, Owner empty
```

## Preconditions

- Catalog: `.tmp` → Tmp|Cache, Owner none.

## Steps

1. Set path under `.tmp`.
2. Expect multi-flag string in ascending bit order; empty owner.

## Context

- Bit order: tmp precedes cache; String is `tmp|cache`.
- Generic home temp dir is not bound to a tool owner.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".tmp/foo.png"
	return nil
}
```
