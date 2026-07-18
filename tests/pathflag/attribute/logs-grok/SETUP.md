# Scenario

**Feature**: `.grok/logs` are Grok application logs

```
caller -> Classify(".grok/logs/app.log")
  -> Rule=.grok/logs, Flags=logs, Owner=grok
```

## Preconditions

- Catalog: `.grok/logs` → Logs, owner grok.

## Steps

1. Set path under `.grok/logs`.
2. Expect logs flag and grok owner.

## Context

- Path rule should win over the generic `**/*.log` suffix for files under this dir.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".grok/logs/app.log"
	return nil
}
```
