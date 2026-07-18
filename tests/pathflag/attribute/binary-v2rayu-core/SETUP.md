# Scenario

**Feature**: only `.V2rayU/v2ray-core` is binary; not whole `.V2rayU`

```
caller -> Classify(".V2rayU/v2ray-core/bin")
  -> Rule=.V2rayU/v2ray-core, Flags=binary, Owner empty
```

## Preconditions

- Catalog: `.V2rayU/v2ray-core` → Binary.
- Config and other files under `.V2rayU` stay unflagged (see no-hit/v2rayu-config).

## Steps

1. Set path under `.V2rayU/v2ray-core`.
2. Expect binary flag only.

## Context

- Reinstallable core binary tree; user config must not inherit this rule.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".V2rayU/v2ray-core/bin"
	return nil
}
```
