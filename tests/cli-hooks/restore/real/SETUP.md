# Scenario

**Feature**: real restore may execute `restoreCmd` (writes allowed)

```
# real restore runs restoreCmd
operator -> bak-files restore --config …
  -> restoreCmd side-effect file appears
```

## Preconditions

- Args without `--dry-run`.

## Steps

1. Leaf seeds config + optional backup content; Assert checks side-effect path.

## Context

- Sibling under restore for future dry-run restoreCmd leaves if added later.

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Real restore: strip accidental --dry-run (leaves set full Args).
	filtered := req.Args[:0]
	for _, a := range req.Args {
		if a == "--dry-run" {
			continue
		}
		filtered = append(filtered, a)
	}
	req.Args = filtered
	_ = strings.Join(req.Args, " ")
	return nil
}
```
