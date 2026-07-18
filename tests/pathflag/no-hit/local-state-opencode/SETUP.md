# Scenario

**Feature**: `.local/state/opencode` is not a fine-prefix skip path

```
caller -> Classify(".local/state/opencode/kv.json")
  -> zero Result (no attribute skip; not under catalogued .local/* rules)
```

## Preconditions

- Fine prefixes only under `.local` (e.g. share/containers, share/claude, Trash) —
  never whole `.local` and not `.local/state/...`.

## Steps

1. Set path to `.local/state/opencode/kv.json`.
2. Expect no skip flags / zero Result.

## Context

- Protects state KV from being swept by a broad `.local` rule.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/state/opencode/kv.json"
	return nil
}
```
