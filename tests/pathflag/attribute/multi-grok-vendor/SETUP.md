# Scenario

**Feature**: `.grok/vendor` is Cache|Vendor for Grok

```
caller -> Classify(".grok/vendor/pkg")
  -> Rule=.grok/vendor, Flags=cache|vendor, Owner=grok
```

## Preconditions

- Catalog: `.grok/vendor` → Cache|Vendor, owner grok.

## Steps

1. Set path under `.grok/vendor`.
2. Expect multi-flag `cache|vendor`.

## Context

- Vendor bit is last in name order; cache comes first.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".grok/vendor/pkg"
	return nil
}
```
