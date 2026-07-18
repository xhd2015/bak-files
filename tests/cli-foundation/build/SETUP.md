# Scenario

**Feature**: module builds the bak-files binary from cmd/bak-files

```
# developer builds from module root
developer -> go build -o <bin> ./cmd/bak-files
  -> binary exists (exit 0)
```

## Preconditions

- Module path `github.com/xhd2015/bak-files`.
- Entry package under `cmd/bak-files`.

## Steps

1. Child sets `Mode` to `build`.
2. `Run` executes a fresh `go build` into a temp path.

## Context

- Independent of session-cached CLI binary used by other leaves (fresh build check).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "build"
	return nil
}
```
