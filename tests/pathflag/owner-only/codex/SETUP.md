# Scenario

**Feature**: `.codex` prefix owns nested non-catalog paths

```
caller -> Classify(".codex/sessions/abc.json")
  -> Owner=codex, Flags empty, Rule empty
```

## Preconditions

- Catalog rules under codex are only `.codex/.tmp` and `.codex/skills/.system`.

## Steps

1. Set path under `.codex/sessions/`.
2. Expect owner-only Result.

## Context

- Distinguishes owner attribution from cache/tmp attribute rows.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".codex/sessions/abc.json"
	return nil
}
```
