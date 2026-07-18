# Scenario

**Feature**: `.android/cache` is Android SDK manager index cache

```
caller -> Classify(".android/cache/sdkbin-1_xxx.xml")
  -> Rule=.android/cache, Flags=cache, Owner=android
```

## Preconditions

- Catalog: `.android/cache` → Cache, Owner android.

## Steps

1. Set path under `.android/cache`.
2. Expect cache flags and owner android.

## Context

- Does not classify `.android/adbkey` (narrower than whole `.android`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".android/cache/sdkbin-1_xxx.xml"
	return nil
}
```
