# Scenario

**Feature**: `.local/share/opencode` owner prefix without attribute hit

```
caller -> Classify(".local/share/opencode/config.json")
  -> Owner=opencode, Flags empty, Rule empty
```

## Preconditions

- Attribute rules exist for `repos`, `snapshot`, `log` under this prefix — not
  for arbitrary files like `config.json`.

## Steps

1. Set path to `.local/share/opencode/config.json`.
2. Expect owner-only Result.

## Context

- Validates longest owner prefix `.local/share/opencode` (not only `.opencode`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/share/opencode/config.json"
	return nil
}
```
