# Scenario

**Feature**: `.codex/.tmp` is Codex temporary plugin cache (multi-flag)

```
caller -> Classify(".codex/.tmp/plugin")
  -> Rule=.codex/.tmp, Flags=tmp|cache, Owner=codex
```

## Preconditions

- Catalog: `.codex/.tmp` → Tmp|Cache, owner codex.

## Steps

1. Set path under `.codex/.tmp`.
2. Expect multi-flag string in ascending bit order.

## Context

- Demonstrates multi-bit Flags formatting (`tmp` before `cache`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".codex/.tmp/plugin"
	return nil
}
```
