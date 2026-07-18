# Scenario

**Feature**: `.local/share/containers` is container/runtime cache (fine prefix)

```
caller -> Classify(".local/share/containers/podman/x")
  -> Rule=.local/share/containers, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.local/share/containers` → Cache, no owner.
- Whole `.local` is **not** a catalog rule (fine prefixes only).

## Steps

1. Set nested path under `.local/share/containers`.
2. Expect cache flag only.

## Context

- Podman/container layers; not history.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/share/containers/podman/x"
	return nil
}
```
