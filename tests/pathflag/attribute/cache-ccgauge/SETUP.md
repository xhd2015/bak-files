# Scenario

**Feature**: `.ccgauge` is ccgauge cache

```
caller -> Classify(".ccgauge/cache/x")
  -> Rule=.ccgauge, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.ccgauge` → Cache, no owner.

## Steps

1. Set nested path under `.ccgauge`.
2. Expect cache flag.

## Context

- Tool-local cache tree under home.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".ccgauge/cache/x"
	return nil
}
```
