# Scenario

**Feature**: any `node_modules` path component matches segment rule

```
caller -> Classify("src/app/node_modules/left-pad/index.js")
  -> Rule=**/node_modules, Flags=vendor, Owner empty
```

## Preconditions

- No path catalog rule applies; segment scan finds `node_modules`.

## Steps

1. Set a deep path with a `node_modules` segment.
2. Expect Vendor via `**/node_modules`.

## Context

- Segment rules apply only when no path catalog rule matched.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "src/app/node_modules/left-pad/index.js"
	return nil
}
```
