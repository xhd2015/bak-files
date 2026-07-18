# Scenario

**Feature**: `.Trash` is macOS trash

```
caller -> Classify(".Trash/file.txt")
  -> Rule=.Trash, Flags=trash, Owner empty
```

## Preconditions

- Catalog: `.Trash` → Trash, Owner none.

## Steps

1. Set path under `.Trash`.
2. Expect trash flag, no owner.

## Context

- Linux trash is a separate rule (`.local/share/Trash`); this leaf is macOS.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".Trash/file.txt"
	return nil
}
```
