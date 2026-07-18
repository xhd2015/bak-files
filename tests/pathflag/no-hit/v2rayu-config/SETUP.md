# Scenario

**Feature**: `.V2rayU` config is not binary/skip from the v2ray-core rule

```
caller -> Classify(".V2rayU/config.json")
  -> zero Result (no Rule/Flags from new rules)
```

## Preconditions

- Only `.V2rayU/v2ray-core` is catalogued as binary — not whole `.V2rayU`.

## Steps

1. Set path to `.V2rayU/config.json`.
2. Expect zero Result (no skip attributes).

## Context

- Negative for over-broad `.V2rayU` prefix matching.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".V2rayU/config.json"
	return nil
}
```
