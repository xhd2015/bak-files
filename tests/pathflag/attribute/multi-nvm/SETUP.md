# Scenario

**Feature**: `.nvm` is Node Version Manager install cache and binaries

```
caller -> Classify(".nvm/versions/node/v20/bin/node")
  -> Rule=.nvm, Flags=cache|binary, Owner=nvm
```

## Preconditions

- Catalog: `.nvm` → Cache|Binary, owner nvm.

## Steps

1. Set nested path under `.nvm`.
2. Expect multi-flag string with cache before binary; owner nvm.

## Context

- Bit order: cache precedes binary; String is `cache|binary`.
- Whole-tree nvm installs are reclaimable/cacheable binaries.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".nvm/versions/node/v20/bin/node"
	return nil
}
```
