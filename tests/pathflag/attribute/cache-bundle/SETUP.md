# Scenario

**Feature**: `.bundle` is Bundler cache

```
caller -> Classify(".bundle/cache/x")
  -> Rule=.bundle, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.bundle` → Cache, no owner.

## Steps

1. Set nested path under `.bundle`.
2. Expect cache flag.

## Context

- Ruby Bundler home cache directory.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".bundle/cache/x"
	return nil
}
```
