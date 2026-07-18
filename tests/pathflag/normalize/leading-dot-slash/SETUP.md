# Scenario

**Feature**: strip leading `./` before classify

```
caller -> Classify("./.bun/install")
  -> Rule=.bun, Flags=cache, Owner=bun
```

## Preconditions

- Leading `./` is removed during normalize.

## Steps

1. Set `RelPath` with leading `./`.
2. Expect same Result as bare `.bun/...`.

## Context

- Common when joining relative fragments.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = "./.bun/install"
	return nil
}
```
