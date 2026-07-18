# Scenario

**Feature**: `.grok/bundled` is reinstallable bundled product cache/binary/vendor

```
caller -> Classify(".grok/bundled/skills/x")
  -> Rule=.grok/bundled, Flags=cache|binary|vendor, Owner=grok
```

## Preconditions

- Catalog: `.grok/bundled` → Cache|Binary|Vendor, Owner grok.

## Steps

1. Set path under `.grok/bundled`.
2. Expect multi flags in bit order and owner grok.

## Context

- Bundled skills/roles ship with the product; skip bak by default.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".grok/bundled/skills/x"
	return nil
}
```
