# Scenario

**Feature**: `.backup` is machine backup metadata

```
caller -> Classify(".backup/index")
  -> Rule=.backup, Flags=meta, Owner empty
```

## Preconditions

- Catalog: `.backup` → Meta, Owner none.

## Steps

1. Set path under `.backup`.
2. Expect meta flag.

## Context

- Meta is for pack-time injected backup metadata, not tool owner data.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".backup/index"
	return nil
}
```
