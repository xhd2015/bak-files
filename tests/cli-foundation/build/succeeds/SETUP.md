# Scenario

**Feature**: `go build ./cmd/bak-files` succeeds

```
# developer builds bak-files
developer -> go build -o <temp>/bak-files-build-check ./cmd/bak-files
  -> exit 0, binary file present
```

## Preconditions

- `Mode=build`.

## Steps

1. Confirm Mode is build (grouping already set).

## Context

- Matches exit criterion: `go build -o /tmp/bak-files ./cmd/bak-files` succeeds.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "build"
	return nil
}
```
