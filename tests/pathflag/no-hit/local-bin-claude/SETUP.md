# Scenario

**Feature**: `.local/bin/claude` is not under the `.local/share/claude` cache rule

```
caller -> Classify(".local/bin/claude")
  -> zero Result (no skip flags from claude share rule)
```

## Preconditions

- Catalog rule is `.local/share/claude` only — must not match `.local/bin/claude`.

## Steps

1. Set path to `.local/bin/claude`.
2. Expect no skip attributes / zero Result.

## Context

- Negative for accidental prefix/substring match on "claude".

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".local/bin/claude"
	return nil
}
```
