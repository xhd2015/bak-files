# Scenario

**Feature**: `.expo` is Expo tooling cache

```
caller -> Classify(".expo/ios-simulator-app-cache/x")
  -> Rule=.expo, Flags=cache, Owner empty
```

## Preconditions

- Catalog: `.expo` → Cache, no owner.

## Steps

1. Set nested path under `.expo`.
2. Expect cache flag.

## Context

- Simulator/app caches under the Expo home dir.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".expo/ios-simulator-app-cache/x"
	return nil
}
```
