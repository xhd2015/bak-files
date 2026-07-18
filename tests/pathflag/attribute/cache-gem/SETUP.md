# Scenario

**Feature**: `.gem` is RubyGems cache

```
caller -> Classify(".gem/specs/x")
  -> Rule=.gem, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.gem` → Cache, no owner.

## Steps

1. Set nested path under `.gem`.
2. Expect cache flag.

## Context

- Specs and installed gem artifacts under home `.gem`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".gem/specs/x"
	return nil
}
```
