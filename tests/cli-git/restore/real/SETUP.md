# Scenario

**Feature**: real restore (no `--dry-run`) may write operator source paths

```
# real restore
operator -> bak-files restore --config …
  -> writes destination from targetDir
```

## Preconditions

- DryRun is false.

## Steps

1. Leaf seeds backup artifact and destination.

## Context

- Real write path only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
