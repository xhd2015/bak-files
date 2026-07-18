# Scenario

**Feature**: cursor-agent versions path is Binary|Cache

```
caller -> Classify(".local/share/cursor-agent/versions/1.0.0")
  -> Rule=.local/share/cursor-agent/versions
  -> Flags=cache|binary, Owner=cursor
```

## Preconditions

- Catalog: `.local/share/cursor-agent/versions` → Binary|Cache, owner cursor.

## Steps

1. Set nested path under versions.
2. Expect multi-flag string with cache before binary.

## Context

- Bit order: cache precedes binary; String is `cache|binary`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/share/cursor-agent/versions/1.0.0"
	return nil
}
```
