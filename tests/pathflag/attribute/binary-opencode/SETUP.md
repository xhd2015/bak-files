# Scenario

**Feature**: `.opencode/bin` is reinstallable OpenCode binary

```
caller -> Classify(".opencode/bin")
  -> Rule=.opencode/bin, Flags=binary, Owner=opencode
```

## Preconditions

- Catalog: `.opencode/bin` → Binary, owner opencode.

## Steps

1. Set exact rule path `.opencode/bin`.
2. Expect binary flag.

## Context

- Exact rule match (not only nested under rule).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".opencode/bin"
	return nil
}
```
