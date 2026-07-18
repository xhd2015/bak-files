# Scenario

**Feature**: `.grok/sessions` are Grok session history (not cache)

```
caller -> Classify(".grok/sessions/%2Ffoo/id")
  -> Rule=.grok/sessions, Flags=history, Owner=grok
```

## Preconditions

- Catalog: `.grok/sessions` → History, owner grok.
- Sessions must not be labeled cache or tmp.

## Steps

1. Set path under `.grok/sessions`.
2. Expect history flag and grok owner.

## Context

- Browser-history-like: keep local, skip bak by default; not reclaim-as-cache.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".grok/sessions/%2Ffoo/id"
	return nil
}
```
