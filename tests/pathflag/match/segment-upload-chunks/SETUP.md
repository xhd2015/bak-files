# Scenario

**Feature**: any `upload-chunks` path component matches segment rule

```
caller -> Classify("var/tmp/upload-chunks/part-001")
  -> Rule=**/upload-chunks, Flags=tmp, Owner empty
```

## Preconditions

- Segment rule after path catalog miss: `upload-chunks` → Tmp.

## Steps

1. Set path containing `upload-chunks`.
2. Expect Tmp via `**/upload-chunks`.

## Context

- Incomplete upload temp state; independent of tool owners.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "var/tmp/upload-chunks/part-001"
	return nil
}
```
